package bulkload

import (
	"encoding/binary"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_RangeReader(t *testing.T) {
	require.NoError(t, fdb.APIVersion(700), "use fdb 7.0")
	db, err := fdb.OpenDefault()
	require.NoError(t, err, "open foundationdb")

	sub := subspace.Sub("Test_RangeReader")

	t.Run("init", func(t *testing.T) {
		for j := 0; j < 20; j++ {
			_, err = db.Transact(func(tr fdb.Transaction) (interface{}, error) {
				for i := uint32(j * 50000); i < uint32(j*50000)+50000; i++ {
					result := make([]byte, 4)
					binary.BigEndian.PutUint32(result, i)
					tr.Set(sub.Sub(result), []byte{})
				}
				return nil, nil
			})
			require.NoError(t, err, "Setup")
		}
	})

	t.Run("test", func(t *testing.T) {
		cnt := 0
		require.NoError(t, New[fdb.KeyValue](db, Reader[fdb.KeyValue](NewRangeReader(db, sub, RangeReaderWithBatchSize(10000))), func(tr fdb.Transaction, value fdb.KeyValue) error {
			cnt++
			return nil
		}).Run(WithBufSize(1000000), WithBatchSize(100000)))
		assert.Equal(t, 1000000, cnt)
	})

	t.Run("cleanup", func(t *testing.T) {
		_, err := db.Transact(func(tr fdb.Transaction) (interface{}, error) {
			tr.Clear(sub)
			return nil, nil
		})
		require.NoError(t, err)
	})

}
