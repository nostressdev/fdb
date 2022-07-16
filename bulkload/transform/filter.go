package transform

import "github.com/nostressdev/fdb/bulkload"

func Filter[T any](producer bulkload.Producer[T], binaryOp func(T) bool) bulkload.Producer[T] {
	return func(ch chan T) error {
		ch1 := make(chan T)
		defer close(ch)
		go func() {
			for task := range ch1 {
				if binaryOp(task) {
					ch <- task
				}
			}
		}()
		return producer(ch1)
	}
}
