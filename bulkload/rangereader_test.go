package bulkload

import (
	"fmt"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
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
		require.NoError(t, New[int](db,
			ListProducer(kTestSize),
			func(tr fdb.Transaction, value int) error {
				tr.Set(sub.Pack(tuple.Tuple{value}), []byte{})
				return nil
			},
		).Run(
			WithConsumersOption(kMaxConsumers),
			WithBatchSize(2500),
			WithBufSize(kTestSize),
		))
	})

	tests := []struct {
		options RangeReaderOptions
	}{
		{options: RangeReaderOptions{batchSize: 1000, producers: 1}},
		{options: RangeReaderOptions{batchSize: 10000000, producers: 1}},
		{options: RangeReaderOptions{batchSize: 1000, producers: 50}},
		{options: RangeReaderOptions{batchSize: 1000000, producers: 50}},
	}

	for index := range tests {
		options := tests[index].options
		t.Run(fmt.Sprintf("Test #%v", index), func(t *testing.T) {
			cnt := 0
			begin, end := sub.FDBRangeKeys()
			require.NoError(t, New[fdb.KeyValue](db, Producer[fdb.KeyValue](
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
			).Run(WithBufSize(kTestSize), WithBatchSize(kTestSize)))
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
