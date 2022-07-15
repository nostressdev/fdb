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

	t.Run("test", func(t *testing.T) {
		type S struct {
			x int
		}
		q := New[*S](sub)
		r, err := q.Dequeue(db)
		require.NoError(t, err, "dequeue empty queue")
		require.Equal(t, nil, r, "equal dequeue empty queue")
		for i := 0; i < 5; i++ {
			require.NoError(t, q.Enqueue(db, &S{x: i}), fmt.Sprintf("enqueue %d", i))
		}
		for i := 5; i > 0; i-- {
			r, err = q.Dequeue(db)
			require.NoError(t, err, fmt.Sprintf("dequeue %d", i))
			require.Equal(t, i, r.x, "equal queue res")
		}
		r, err = q.Dequeue(db)
		require.NoError(t, err, "dequeue empty queue")
		require.Equal(t, nil, r, "equal dequeue empty queue")
	})

	t.Run("cleanup", func(t *testing.T) {
		_, err := db.Transact(func(tr fdb.Transaction) (interface{}, error) {
			tr.Clear(sub)
			return nil, nil
		})
		require.NoError(t, err)
	})

}
