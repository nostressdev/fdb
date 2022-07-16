package future

import "github.com/nostressdev/fdb/utils/errors"

type Future[Output any] struct {
	result chan futureResult[Output]
}

func (future *Future[Output]) Close() {
	close(future.result)
}

func (future *Future[Output]) Await() (Output, error) {
	select {
	case result, ok := <-future.result:
		if !ok {
			return result.value, errors.RuntimeError.New("future was awaited twice")
		}
		return result.value, result.err
	}
}

func AwaitAll[Output any](futures ...Future[Output]) ([]Output, error) {
	result := make([]Output, 0, len(futures))
	for _, future := range futures {
		if value, err := future.Await(); err != nil {
			return nil, err
		} else {
			result = append(result, value)
		}
	}
	return result, nil
}

type futureResult[T any] struct {
	value T
	err   error
}

func Async[T any](f func() (T, error)) Future[T] {
	future := Future[T]{
		result: make(chan futureResult[T]),
	}
	go func() {
		value, err := f()
		future.result <- futureResult[T]{
			value: value,
			err:   err,
		}
		future.Close()
	}()
	return future
}

type FutureNil struct {
	future Future[Nil]
}

func (future *FutureNil) Await() error {
	_, err := future.future.Await()
	return err
}

func AwaitAllNil(futures ...FutureNil) error {
	result := make([]Future[Nil], 0, len(futures))
	for _, future := range futures {
		result = append(result, future.future)
	}
	_, err := AwaitAll(result...)
	return err
}

func AsyncNil(f func() error) FutureNil {
	return FutureNil{
		future: Async[Nil](func() (Nil, error) {
			return new(Nil), f()
		}),
	}
}

type Nil interface{}
