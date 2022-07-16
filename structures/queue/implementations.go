package queue

import (
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
)

func NewInt(sub subspace.Subspace) Queue[int] {
	return Queue[int]{
		sub: sub,
		pack: func(i int) ([]byte, error) {
			return tuple.Tuple{i}.Pack(), nil
		},
		unpack: func(v []byte) (int, error) {
			res, err := tuple.Unpack(v)
			if err != nil {
				return 0, err
			}
			return int(res[0].(int64)), nil
		},
	}
}

func NewInt64(sub subspace.Subspace) Queue[int64] {
	return Queue[int64]{
		sub: sub,
		pack: func(i int64) ([]byte, error) {
			return tuple.Tuple{i}.Pack(), nil
		},
		unpack: func(v []byte) (int64, error) {
			res, err := tuple.Unpack(v)
			if err != nil {
				return 0, err
			}
			return res[0].(int64), nil
		},
	}
}

func NewUint(sub subspace.Subspace) Queue[uint] {
	return Queue[uint]{
		sub: sub,
		pack: func(i uint) ([]byte, error) {
			return tuple.Tuple{i}.Pack(), nil
		},
		unpack: func(v []byte) (uint, error) {
			res, err := tuple.Unpack(v)
			if err != nil {
				return 0, err
			}
			return uint(res[0].(int64)), nil
		},
	}
}

func NewUint64(sub subspace.Subspace) Queue[uint64] {
	return Queue[uint64]{
		sub: sub,
		pack: func(i uint64) ([]byte, error) {
			return tuple.Tuple{i}.Pack(), nil
		},
		unpack: func(v []byte) (uint64, error) {
			res, err := tuple.Unpack(v)
			if err != nil {
				return 0, err
			}
			if v, ok := res[0].(int64); ok {
				return uint64(v), nil
			}
			return res[0].(uint64), nil
		},
	}
}

func NewString(sub subspace.Subspace) Queue[string] {
	return Queue[string]{
		sub: sub,
		pack: func(i string) ([]byte, error) {
			return []byte(i), nil
		},
		unpack: func(v []byte) (string, error) {
			return string(v), nil
		},
	}
}

func NewBytes(sub subspace.Subspace) Queue[[]byte] {
	return Queue[[]byte]{
		sub: sub,
		pack: func(i []byte) ([]byte, error) {
			return i, nil
		},
		unpack: func(v []byte) ([]byte, error) {
			return v, nil
		},
	}
}
