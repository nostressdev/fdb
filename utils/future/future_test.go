package future

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Future(t *testing.T) {
	expectedErr := errors.New("test error")
	f1 := Async[int](func() (int, error) {
		return 1, nil
	})
	f2 := Async[int](func() (int, error) {
		return 2, nil
	})
	f3 := Async[int](func() (int, error) {
		return 3, expectedErr
	})
	value, err := f2.Await()
	assert.EqualValues(t, 2, value)
	assert.NoError(t, err)
	value, err = f1.Await()
	assert.EqualValues(t, 1, value)
	assert.NoError(t, err)
	value, err = f3.Await()
	assert.EqualValues(t, 3, value)
	assert.Error(t, err)
}

func Test_FutureNil(t *testing.T) {
	expectedErr := errors.New("test error")
	f1 := AsyncNil(func() error {
		return nil
	})
	f2 := AsyncNil(func() error {
		return expectedErr
	})
	err := f1.Await()
	assert.NoError(t, err)
	err = f2.Await()
	assert.Error(t, err)
}

func Test_FutureAwaitAll(t *testing.T) {
	expectedErr := errors.New("test error")
	f1 := Async[int](func() (int, error) {
		return 1, nil
	})
	f2 := Async[int](func() (int, error) {
		return 2, nil
	})
	f3 := Async[int](func() (int, error) {
		return 3, expectedErr
	})
	_, err := AwaitAll(f1, f2, f3)
	assert.Error(t, err)
}

func Test_FutureAwaitAllNil(t *testing.T) {
	expectedErr := errors.New("test error")
	f1 := AsyncNil(func() error {
		return nil
	})
	f2 := AsyncNil(func() error {
		return nil
	})
	f3 := AsyncNil(func() error {
		return expectedErr
	})
	assert.Error(t, AwaitAllNil(f1, f2, f3))
}

func Test_FutureAwaitAllValues(t *testing.T) {
	f1 := Async[int](func() (int, error) {
		return 1, nil
	})
	f2 := Async[int](func() (int, error) {
		return 2, nil
	})
	f3 := Async[int](func() (int, error) {
		return 3, nil
	})
	value, err := AwaitAll(f1, f2, f3)
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, value)
}

func Test_FutureAwaitTwice(t *testing.T) {
	f := AsyncNil(func() error {
		return nil
	})
	assert.NoError(t, f.Await())
	assert.Error(t, f.Await())
}
