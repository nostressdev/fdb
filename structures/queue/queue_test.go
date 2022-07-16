package queue

import (
	"fmt"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/stretchr/testify/require"
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

	t.Run("test int64", func(t *testing.T) {
		q := New[int64](sub)
		r, err := q.Dequeue(db)
		require.NoError(t, err, "dequeue empty queue")
		require.Equal(t, true, r == nil, "equal dequeue empty queue")
		for i := int64(0); i < 5; i++ {
			require.NoError(t, q.Enqueue(db, i), fmt.Sprintf("enqueue %d", i))
		}
		for i := int64(0); i < 5; i++ {
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

	t.Run("test string", func(t *testing.T) {
		q := New[string](sub)
		str := []string{"a", "b", "c", "d", "e"}
		r, err := q.Dequeue(db)
		require.NoError(t, err, "dequeue empty queue")
		require.Equal(t, true, r == nil, "equal dequeue empty queue")
		for i := 0; i < 5; i++ {
			require.NoError(t, q.Enqueue(db, str[i]), fmt.Sprintf("enqueue %d", i))
		}
		for i := 0; i < 5; i++ {
			r, err = q.Dequeue(db)
			require.NoError(t, err, fmt.Sprintf("dequeue %d", i))
			require.Equal(t, str[i], *r, "equal queue res")
		}
		r, err = q.Dequeue(db)
		require.NoError(t, err, "dequeue empty queue")
		require.Equal(t, true, r == nil, "equal dequeue empty queue")
	})

	t.Run("test int", func(t *testing.T) {
		q := New[int](sub)
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

	t.Run("test keyConvertible", func(t *testing.T) {
		values := []fdb.KeyConvertible{fdb.Key([]byte{1}), fdb.Key([]byte{2}), fdb.Key([]byte{3}), fdb.Key([]byte{4}), fdb.Key([]byte{5})}
		q := New[fdb.KeyConvertible](sub)
		r, err := q.Dequeue(db)
		require.NoError(t, err, "dequeue empty queue")
		require.Equal(t, true, r == nil, "equal dequeue empty queue")
		for i := 0; i < 5; i++ {
			require.NoError(t, q.Enqueue(db, values[i]), fmt.Sprintf("enqueue %d", i))
		}
		for i := 0; i < 5; i++ {
			r, err = q.Dequeue(db)
			require.NoError(t, err, fmt.Sprintf("dequeue %d", i))
			require.Equal(t, values[i], *r, "equal queue res")
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
