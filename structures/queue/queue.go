package queue

import (
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"math/rand"
)

type Queue[T any] struct {
	sub    subspace.Subspace
	pack   func(T) ([]byte, error)
	unpack func([]byte) (T, error)
}

func New[T any](sub subspace.Subspace, pack func(T) ([]byte, error), unpack func([]byte) (T, error)) Queue[T] {
	return Queue[T]{
		sub:    sub,
		pack:   pack,
		unpack: unpack,
	}
}

func (q *Queue[T]) Dequeue(transactor fdb.Transactor) (*T, error) {
	res, err := transactor.Transact(func(tr fdb.Transaction) (interface{}, error) {
		kv, err := q.firstItem(tr)
		if err != nil {
			return nil, err
		}
		if kv == nil {
			return nil, nil
		}
		tr.Clear(kv.Key)
		res, err := q.unpack(kv.Value)
		if err != nil {
			return nil, err
		}
		return res, nil
	})
	if res != nil {
		resT := res.(T)
		return &resT, err
	}
	return nil, err
}

func (q *Queue[T]) Enqueue(transactor fdb.Transactor, t T) error {
	_, err := transactor.Transact(func(tr fdb.Transaction) (interface{}, error) {
		i, err := q.lastIndex(tr)
		if err != nil {
			return nil, err
		}
		bytes := make([]byte, 20)
		rand.Read(bytes)
		v, err := q.pack(t)
		if err != nil {
			return nil, err
		}
		tr.Set(q.sub.Sub(i+1, bytes), v)
		return nil, nil
	})
	return err
}

func (q *Queue[T]) firstItem(transactor fdb.ReadTransactor) (*fdb.KeyValue, error) {
	res, err := transactor.ReadTransact(func(tr fdb.ReadTransaction) (interface{}, error) {
		iter := tr.GetRange(q.sub, fdb.RangeOptions{
			Mode:  fdb.StreamingModeWantAll,
			Limit: 1,
		}).Iterator()
		if iter.Advance() {
			return iter.Get()
		}
		return nil, nil
	})
	if res != nil {
		resKV := res.(fdb.KeyValue)
		return &resKV, err
	}
	return nil, err
}

func (q *Queue[T]) lastIndex(transactor fdb.ReadTransactor) (int64, error) {
	res, err := transactor.ReadTransact(func(tr fdb.ReadTransaction) (interface{}, error) {
		iter := tr.GetRange(q.sub, fdb.RangeOptions{
			Mode:    fdb.StreamingModeWantAll,
			Limit:   1,
			Reverse: true,
		}).Iterator()
		if iter.Advance() {
			kv, err := iter.Get()
			if err != nil {
				return nil, err
			}
			t, err := q.sub.Unpack(kv.Key)
			if err != nil {
				return nil, err
			}
			return t[0].(int64), nil
		}
		return nil, nil
	})
	if res != nil {
		return res.(int64), err
	}
	return 0, err
}
