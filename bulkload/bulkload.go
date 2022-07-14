package bulkload

import (
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"sync"
)

type BulkLoad[T any] interface {
	Run(opts ...Options) error
}

type Reader[T any] func(chan T) error
type Writer[T any] func(transaction fdb.Transaction, value T) error

type bulkLoadImpl[T any] struct {
	reader Reader[T]
	writer Writer[T]
	db     fdb.Database
}

func New[T any](db fdb.Database, reader Reader[T], writer Writer[T]) BulkLoad[T] {
	return &bulkLoadImpl[T]{
		reader: reader,
		writer: writer,
		db:     db,
	}
}

func (bl *bulkLoadImpl[T]) Run(opts ...Options) error {
	options := mergeOptions(opts...)
	tasks := make(chan T, options.bufSize)
	errCh := make(chan error)

	go func() {
		err := bl.reader(tasks)
		close(tasks)
		if err != nil {
			errCh <- err
		}
	}()

	wgConsumers := new(sync.WaitGroup)
	wgConsumers.Add(options.consumers)
	for consumerIdx := 0; consumerIdx < options.consumers; consumerIdx++ {
		go func() {
			taskIter := make([]T, 0, options.batchSize)
			processTasks := func(taskIter []T) error {
				if len(taskIter) == 0 {
					return nil
				}
				_, err := bl.db.Transact(func(transaction fdb.Transaction) (interface{}, error) {
					for _, task := range taskIter {
						if err := bl.writer(transaction, task); err != nil {
							return nil, err
						}
					}
					return nil, nil
				})
				return err
			}
			for task := range tasks {
				taskIter = append(taskIter, task)
				if len(taskIter) == options.batchSize {
					if err := processTasks(taskIter); err != nil {
						wgConsumers.Done()
						errCh <- err
						return
					}
					taskIter = nil
				}
			}
			if err := processTasks(taskIter); err != nil {
				wgConsumers.Done()
				errCh <- err
				return
			}
			wgConsumers.Done()
		}()
	}
	wgConsumers.Wait()
	close(errCh)
	return <-errCh
}
