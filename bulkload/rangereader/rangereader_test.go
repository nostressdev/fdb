package rangereader

import (
	"fmt"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
	"github.com/nostressdev/fdb/bulkload"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_RangeReader(t *testing.T) {
	const kTestSize = 1000000
	const kMaxConsumers = 50

	require.NoError(t, fdb.APIVersion(700), "use fdb 7.0")
	db, err := fdb.OpenDefault()
	require.NoError(t, err, "open foundationdb")

	sub := subspace.Sub("Test_RangeReader")

	t.Run("init", func(t *testing.T) {
		require.NoError(t, bulkload.New[int](db,
			bulkload.ListProducer(kTestSize),
			func(tr fdb.Transaction, value int) error {
				tr.Set(sub.Pack(tuple.Tuple{value}), []byte{})
				return nil
			},
		).Run(
			bulkload.WithConsumersOption(kMaxConsumers),
			bulkload.WithBatchSize(2500),
			bulkload.WithBufSize(kTestSize),
		))
	})

	tests := []struct {
		options Options
	}{
		{options: Options{batchSize: 1000, producers: 1}},
		{options: Options{batchSize: 10000000, producers: 1}},
		{options: Options{batchSize: 1000, producers: 50}},
		{options: Options{batchSize: 1000000, producers: 50}},
	}

	for index := range tests {
		options := tests[index].options
		t.Run(fmt.Sprintf("Test #%v", index), func(t *testing.T) {
			cnt := 0
			begin, end := sub.FDBRangeKeys()
			require.NoError(t, bulkload.New[fdb.KeyValue](db, bulkload.Producer[fdb.KeyValue](
				NewRangeReader(db,
					fdb.KeyRange{
						Begin: begin,
						End:   end,
					},
					options,
				)),
				func(tr fdb.Transaction, value fdb.KeyValue) error {
					cnt++
					return nil
				},
			).Run(bulkload.WithBufSize(kTestSize), bulkload.WithBatchSize(kTestSize)))
			assert.Equal(t, kTestSize, cnt)
		})
	}

	t.Run("cleanup", func(t *testing.T) {
		_, err := db.Transact(func(tr fdb.Transaction) (interface{}, error) {
			tr.ClearRange(sub)
			return nil, nil
		})
		require.NoError(t, err)
	})

}
