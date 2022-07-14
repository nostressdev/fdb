package bulkload

import (
	"encoding/binary"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_SimpleWrite(t *testing.T) {
	require.NoError(t, fdb.APIVersion(700), "use fdb 7.0")
	db, err := fdb.OpenDefault()
	require.NoError(t, err, "open foundationdb")
	sub := subspace.Sub("TestSimpleWrite")
	tests := []struct {
		opts []Options
		name string
		cnt  uint32
	}{
		{
			name: "Simple",
			opts: []Options{
				WithBatchSize(2),
			},
			cnt: 5,
		},
		{
			name: "Multiple consumers",
			opts: []Options{
				WithConsumersOption(10),
			},
			cnt: 100,
		},
		{
			name: "Multiple consumers and buf size",
			opts: []Options{
				WithConsumersOption(50),
				WithBatchSize(1000),
				WithBufSize(1000),
			},
			cnt: 1000000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = New[uint32](db, func(ch chan uint32) error {
				for i := uint32(0); i < tt.cnt; i++ {
					ch <- i
				}
				return nil
			}, func(tr fdb.Transaction, value uint32) error {
				result := make([]byte, 4)
				binary.BigEndian.PutUint32(result, value)
				tr.Set(sub.Sub(result), []byte{})
				return nil
			}).Run(tt.opts...)
			require.NoError(t, err, "Run should not fail")
			_, err = db.Transact(func(tr fdb.Transaction) (interface{}, error) {
				iter := tr.GetRange(sub, fdb.RangeOptions{
					Mode: fdb.StreamingModeWantAll,
				}).Iterator()
				idx := uint32(0)
				for iter.Advance() {
					value, err := iter.Get()
					require.NoError(t, err)
					result := make([]byte, 4)
					binary.BigEndian.PutUint32(result, idx)
					assert.Equal(t, sub.Sub(result).FDBKey(), value.Key.FDBKey())
					idx += 1
				}
				assert.Equal(t, idx, tt.cnt)
				return nil, nil
			})
			assert.NoError(t, err, "Read check transaction should not fail")
			_, err = db.Transact(func(tr fdb.Transaction) (interface{}, error) {
				tr.ClearRange(sub)
				return nil, nil
			})
			require.NoError(t, err, "Cleanup")
		})
	}

}
