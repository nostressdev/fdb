package transform

import (
	"github.com/nostressdev/fdb/bulkload"
	"github.com/nostressdev/fdb/utils/future"
)

func Map[T any, T1 any](producer bulkload.Producer[T], mapOp func(T) (T1, error)) bulkload.Producer[T1] {
	return func(ch chan T1) error {
		ch1 := make(chan T)
		defer close(ch)
		mapFuture := future.AsyncNil(func() error {
			for task := range ch1 {
				value, err := mapOp(task)
				if err != nil {
					return err
				}
				ch <- value
			}
			return nil
		})
		if err := producer(ch1); err != nil {
			return err
		}
		return mapFuture.Await()
	}
}
