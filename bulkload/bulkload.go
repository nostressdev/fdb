package bulkload

import (
	"fmt"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/nostressdev/fdb/utils/future"
	"golang.org/x/xerrors"
)

type BulkLoad[T any] interface {
	Run(opts ...Options) error
}

type Producer[T any] func(chan T) error
type Consumer[T any] func(transaction fdb.Transaction, value T) error

type bulkLoadImpl[T any] struct {
	producer Producer[T]
	consumer Consumer[T]
	db       fdb.Database
}

func New[T any](db fdb.Database, producer Producer[T], consumer Consumer[T]) BulkLoad[T] {
	return &bulkLoadImpl[T]{
		producer: producer,
		consumer: consumer,
		db:       db,
	}
}

func (bl *bulkLoadImpl[T]) processTasksSet(tr fdb.Transaction, tasksSet []T, options Options) ([]T, error) {
	maxTasksInBatch := len(tasksSet)
	for {
		tr.Reset()
		var err error
		for index := 0; index < maxTasksInBatch; index++ {
			err = bl.consumer(tr, tasksSet[index])
			if err != nil {
				fmt.Println(err)
				break
			}
		}
		if err == nil {
			err = tr.Commit().Get()
		}
		var trErr fdb.Error
		if err != nil {
			if xerrors.As(err, &trErr) {
				if (trErr.Code >= 2101 && trErr.Code <= 2103) || trErr.Code == 1007 {
					maxTasksInBatch = int(float64(maxTasksInBatch) * options.degradationFactor)
					if maxTasksInBatch == 0 {
						return tasksSet, err
					}
					fmt.Println(trErr.Error())
					continue
				}
				if err = tr.OnError(trErr).Get(); err != nil {
					return tasksSet, err
				}
				fmt.Println(trErr.Error())
				continue
			}
			return tasksSet, err
		}
		return tasksSet[maxTasksInBatch:], nil
	}
}

func (bl *bulkLoadImpl[T]) Run(opts ...Options) error {
	options := mergeOptions(opts...)
	tasks := make(chan T, options.bufSize)
	readTask := future.AsyncNil(func() error {
		defer close(tasks)
		return bl.producer(tasks)
	})
	futures := make([]future.FutureNil, 0, options.consumers)
	for consumerIdx := 0; consumerIdx < options.consumers; consumerIdx++ {
		tr, err := bl.db.CreateTransaction()
		if err != nil {
			return err
		}
		futures = append(futures, future.AsyncNil(func() error {
			tasksSet := make([]T, 0, options.batchSize)
			for task := range tasks {
				tasksSet = append(tasksSet, task)
				if len(tasksSet) == options.batchSize {
					restTasksSet, err := bl.processTasksSet(tr, tasksSet, options)
					if err != nil {
						return err
					}
					tasksSet = restTasksSet
				}
			}
			for len(tasksSet) > 0 {
				restTasksSet, err := bl.processTasksSet(tr, tasksSet, options)
				if err != nil {
					return err
				}
				tasksSet = restTasksSet
			}
			return nil
		}))
	}
	readAllTasksIfErrorFuture := future.AsyncNil(func() error {
		err := future.AwaitAllNil(futures...)
		if err != nil {
			for range tasks {
			}
			return err
		}
		return nil
	})
	err := readTask.Await()
	if err != nil {
		return err
	}
	return readAllTasksIfErrorFuture.Await()
}
