package bulkload

import (
	"encoding/binary"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
	"github.com/nostressdev/fdb/bulkload/rangereader"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_BulkLoadWrite(t *testing.T) {
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
		{
			name: "Transaction limit fails",
			opts: []Options{
				WithConsumersOption(1),
				WithBatchSize(100000),
				WithBufSize(1000000),
				WithDegradationFactor(0.75),
			},
			cnt: 1000000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, New(db, ListProducer(tt.cnt),
				func(tr fdb.Transaction, value uint32) error {
					result := make([]byte, 4)
					binary.BigEndian.PutUint32(result, value)
					tr.Set(sub.Sub(result), result)
					return nil
				},
			).Run(tt.opts...), "Run should not fail")

			subBegin, subEnd := sub.FDBRangeKeys()

			m := map[uint32]struct{}{}

			require.NoError(t, New[fdb.KeyValue](db, Producer[fdb.KeyValue](rangereader.NewRangeReader(db, fdb.KeyRange{
				Begin: subBegin,
				End:   subEnd,
			}, rangereader.WithBatchSize(int(tt.cnt)), rangereader.WithProducersOption(50))), func(tr fdb.Transaction, value fdb.KeyValue) error {
				x := binary.BigEndian.Uint32(value.Value)
				if x < tt.cnt {
					m[x] = struct{}{}
				}
				return nil
			}).Run(WithBatchSize(1000), WithBufSize(int(tt.cnt))))

			assert.Equal(t, int(tt.cnt), len(m), "Every element should be used just once")
			assert.NoError(t, err, "Read check transaction should not fail")
			_, err = db.Transact(func(tr fdb.Transaction) (interface{}, error) {
				tr.ClearRange(sub)
				return nil, nil
			})
			require.NoError(t, err, "Cleanup")
		})
	}

}

func Test_BulkLoadConsumerError(t *testing.T) {
	require.NoError(t, fdb.APIVersion(700), "use fdb 7.0")
	db, err := fdb.OpenDefault()
	require.NoError(t, err, "open foundationdb")

	tests := []struct {
		opts []Options
		name string
		cnt  uint32
	}{
		{
			name: "One consumer",
			opts: []Options{
				WithConsumersOption(1),
			},
			cnt: 50,
		},
		{
			name: "Multiple consumers",
			opts: []Options{
				WithConsumersOption(10),
			},
			cnt: 50,
		},
		{
			name: "Multiple consumers big array",
			opts: []Options{
				WithConsumersOption(10),
			},
			cnt: 500000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Error(t, New(db, ListProducer(100),
				func(tr fdb.Transaction, value int) error {
					return errors.New("test error")
				},
			).Run(tt.opts...), "Run should fail")
		})
	}
}

func Test_BulkLoadProducerError(t *testing.T) {
	require.NoError(t, fdb.APIVersion(700), "use fdb 7.0")
	db, err := fdb.OpenDefault()
	require.NoError(t, err, "open foundationdb")
	tests := []struct {
		opts []Options
		name string
		cnt  uint32
	}{
		{
			name: "One consumer",
			opts: []Options{
				WithConsumersOption(1),
			},
			cnt: 50,
		},
		{
			name: "Multiple consumers",
			opts: []Options{
				WithConsumersOption(10),
			},
			cnt: 50,
		},
		{
			name: "Multiple consumers big array",
			opts: []Options{
				WithConsumersOption(10),
			},
			cnt: 500000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Error(t,
				New(db, func(chan int) error {
					return errors.New("test error")
				},
					func(tr fdb.Transaction, value int) error {
						return nil
					},
				).Run(tt.opts...), "Run should fail")
		})
	}
}

func Test_BulkLoadBothError(t *testing.T) {
	require.NoError(t, fdb.APIVersion(700), "use fdb 7.0")
	db, err := fdb.OpenDefault()
	require.NoError(t, err, "open foundationdb")

	tests := []struct {
		opts []Options
		name string
		cnt  uint32
	}{
		{
			name: "One consumer",
			opts: []Options{
				WithConsumersOption(1),
			},
			cnt: 50,
		},
		{
			name: "Multiple consumers",
			opts: []Options{
				WithConsumersOption(10),
			},
			cnt: 50,
		},
		{
			name: "Multiple consumers big array",
			opts: []Options{
				WithConsumersOption(10),
			},
			cnt: 500000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Error(t,
				New(db, func(chan int) error {
					return errors.New("test error")
				},
					func(tr fdb.Transaction, value int) error {
						return errors.New("test error")
					},
				).Run(tt.opts...), "Run should fail")
		})
	}
}

func Test_OneBigTransaction(t *testing.T) {
	require.NoError(t, fdb.APIVersion(700), "use fdb 7.0")
	db, err := fdb.OpenDefault()
	require.NoError(t, err, "open foundationdb")
	sub := subspace.Sub("TestSimpleWrite")
	require.Error(t,
		New(db,
			func(ch chan int) error {
				ch <- 0
				return nil
			},
			func(tr fdb.Transaction, value int) error {
				for i := 0; i < 1000000; i++ {
					tr.Set(sub.Sub(tuple.Tuple{i}.Pack()), []byte{})
				}
				return nil
			},
		).Run(),
		"Run should fail",
	)
}
