package queue

import (
	"encoding/json"
	"fmt"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func Test_Queue(t *testing.T) {
	require.NoError(t, fdb.APIVersion(700), "use fdb 7.0")
	db, err := fdb.OpenDefault()
	require.NoError(t, err, "open foundationdb")

	sub := subspace.Sub("Test_Queue")

	t.Run("cleanup", func(t *testing.T) {
		_, err := db.Transact(func(tr fdb.Transaction) (interface{}, error) {
			tr.ClearRange(sub)
			return nil, nil
		})
		require.NoError(t, err)
	})

	t.Run("test int", func(t *testing.T) {
		q := NewInt(sub)
		r, err := q.Dequeue(db)
		require.NoError(t, err, "dequeue empty queue")
		require.Equal(t, true, r == nil, "equal dequeue empty queue")
		for i := 0; i < 5; i++ {
			require.NoError(t, q.Enqueue(db, i), fmt.Sprintf("enqueue %d", i))
		}
		for i := 0; i < 5; i++ {
			r, err = q.Dequeue(db)
			require.NoError(t, err, fmt.Sprintf("dequeue %d", i))
			require.Equal(t, i, *r, "equal queue res")
		}
		r, err = q.Dequeue(db)
		require.NoError(t, err, "dequeue empty queue")
		require.Equal(t, true, r == nil, "equal dequeue empty queue")
	})

	t.Run("cleanup", func(t *testing.T) {
		_, err := db.Transact(func(tr fdb.Transaction) (interface{}, error) {
			tr.ClearRange(sub)
			return nil, nil
		})
		require.NoError(t, err)
	})

	t.Run("test int64", func(t *testing.T) {
		q := NewInt64(sub)
		r, err := q.Dequeue(db)
		require.NoError(t, err, "dequeue empty queue")
		require.Equal(t, true, r == nil, "equal dequeue empty queue")
		for i := 0; i < 5; i++ {
			require.NoError(t, q.Enqueue(db, int64(i)), fmt.Sprintf("enqueue %d", i))
		}
		for i := 0; i < 5; i++ {
			r, err = q.Dequeue(db)
			require.NoError(t, err, fmt.Sprintf("dequeue %d", i))
			require.Equal(t, int64(i), *r, "equal queue res")
		}
		r, err = q.Dequeue(db)
		require.NoError(t, err, "dequeue empty queue")
		require.Equal(t, true, r == nil, "equal dequeue empty queue")
	})

	t.Run("cleanup", func(t *testing.T) {
		_, err := db.Transact(func(tr fdb.Transaction) (interface{}, error) {
			tr.ClearRange(sub)
			return nil, nil
		})
		require.NoError(t, err)
	})

	t.Run("test uint", func(t *testing.T) {
		q := NewUint(sub)
		r, err := q.Dequeue(db)
		require.NoError(t, err, "dequeue empty queue")
		require.Equal(t, true, r == nil, "equal dequeue empty queue")
		for i := 0; i < 5; i++ {
			require.NoError(t, q.Enqueue(db, uint(i)), fmt.Sprintf("enqueue %d", i))
		}
		for i := 0; i < 5; i++ {
			r, err = q.Dequeue(db)
			require.NoError(t, err, fmt.Sprintf("dequeue %d", i))
			require.Equal(t, uint(i), *r, "equal queue res")
		}
		r, err = q.Dequeue(db)
		require.NoError(t, err, "dequeue empty queue")
		require.Equal(t, true, r == nil, "equal dequeue empty queue")
	})

	t.Run("cleanup", func(t *testing.T) {
		_, err := db.Transact(func(tr fdb.Transaction) (interface{}, error) {
			tr.ClearRange(sub)
			return nil, nil
		})
		require.NoError(t, err)
	})

	t.Run("test uint64", func(t *testing.T) {
		q := NewUint64(sub)
		r, err := q.Dequeue(db)
		require.NoError(t, err, "dequeue empty queue")
		require.Equal(t, true, r == nil, "equal dequeue empty queue")
		for i := 0; i < 5; i++ {
			require.NoError(t, q.Enqueue(db, uint64(i)), fmt.Sprintf("enqueue %d", i))
		}
		for i := 0; i < 5; i++ {
			r, err = q.Dequeue(db)
			require.NoError(t, err, fmt.Sprintf("dequeue %d", i))
			require.Equal(t, uint64(i), *r, "equal queue res")
		}
		r, err = q.Dequeue(db)
		require.NoError(t, err, "dequeue empty queue")
		require.Equal(t, true, r == nil, "equal dequeue empty queue")
	})

	t.Run("cleanup", func(t *testing.T) {
		_, err := db.Transact(func(tr fdb.Transaction) (interface{}, error) {
			tr.ClearRange(sub)
			return nil, nil
		})
		require.NoError(t, err)
	})

	t.Run("test string", func(t *testing.T) {
		q := NewString(sub)
		r, err := q.Dequeue(db)
		require.NoError(t, err, "dequeue empty queue")
		require.Equal(t, true, r == nil, "equal dequeue empty queue")
		for i := 0; i < 5; i++ {
			require.NoError(t, q.Enqueue(db, strconv.Itoa(i)), fmt.Sprintf("enqueue %d", i))
		}
		for i := 0; i < 5; i++ {
			r, err = q.Dequeue(db)
			require.NoError(t, err, fmt.Sprintf("dequeue %d", i))
			require.Equal(t, strconv.Itoa(i), *r, "equal queue res")
		}
		r, err = q.Dequeue(db)
		require.NoError(t, err, "dequeue empty queue")
		require.Equal(t, true, r == nil, "equal dequeue empty queue")
	})

	t.Run("cleanup", func(t *testing.T) {
		_, err := db.Transact(func(tr fdb.Transaction) (interface{}, error) {
			tr.ClearRange(sub)
			return nil, nil
		})
		require.NoError(t, err)
	})

	t.Run("test []byte", func(t *testing.T) {
		q := NewBytes(sub)
		r, err := q.Dequeue(db)
		require.NoError(t, err, "dequeue empty queue")
		require.Equal(t, true, r == nil, "equal dequeue empty queue")
		for i := 0; i < 5; i++ {
			require.NoError(t, q.Enqueue(db, []byte{byte(i)}), fmt.Sprintf("enqueue %d", i))
		}
		for i := 0; i < 5; i++ {
			r, err = q.Dequeue(db)
			require.NoError(t, err, fmt.Sprintf("dequeue %d", i))
			require.Equal(t, []byte{byte(i)}, *r, "equal queue res")
		}
		r, err = q.Dequeue(db)
		require.NoError(t, err, "dequeue empty queue")
		require.Equal(t, true, r == nil, "equal dequeue empty queue")
	})

	t.Run("cleanup", func(t *testing.T) {
		_, err := db.Transact(func(tr fdb.Transaction) (interface{}, error) {
			tr.ClearRange(sub)
			return nil, nil
		})
		require.NoError(t, err)
	})

	t.Run("test my queue", func(t *testing.T) {
		type S struct {
			I int `json:"i"`
			J int `json:"j"`
		}
		q := New[S](sub,
			func(t S) ([]byte, error) {
				return json.Marshal(t)
			},
			func(v []byte) (S, error) {
				s := S{}
				return s, json.Unmarshal(v, &s)
			})
		r, err := q.Dequeue(db)
		require.NoError(t, err, "dequeue empty queue")
		require.Equal(t, true, r == nil, "equal dequeue empty queue")
		for i := 0; i < 5; i++ {
			require.NoError(t, q.Enqueue(db, S{i, 5 - i}), fmt.Sprintf("enqueue %d", i))
		}
		for i := 0; i < 5; i++ {
			r, err = q.Dequeue(db)
			require.NoError(t, err, fmt.Sprintf("dequeue %d", i))
			require.Equal(t, S{i, 5 - i}, *r, "equal queue res")
		}
		r, err = q.Dequeue(db)
		require.NoError(t, err, "dequeue empty queue")
		require.Equal(t, true, r == nil, "equal dequeue empty queue")
	})

	t.Run("cleanup", func(t *testing.T) {
		_, err := db.Transact(func(tr fdb.Transaction) (interface{}, error) {
			tr.ClearRange(sub)
			return nil, nil
		})
		require.NoError(t, err)
	})
}
