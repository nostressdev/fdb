package bulkload

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MergeOptions(t *testing.T) {
	tests := []struct {
		name   string
		opts   []Options
		result Options
	}{
		{
			name: "Simple case",
			opts: []Options{
				WithConsumersOption(3),
				WithBufSize(123),
				WithBatchSize(11),
				WithDegradationFactor(0.9),
			},
			result: Options{
				consumers:         3,
				bufSize:           123,
				batchSize:         11,
				degradationFactor: 0.9,
			},
		},
		{
			name: "Simple case reversed",
			opts: []Options{
				WithDegradationFactor(0.1),
				WithBufSize(123),
				WithConsumersOption(3),
				WithBatchSize(11),
			},
			result: Options{
				bufSize:           123,
				consumers:         3,
				batchSize:         11,
				degradationFactor: 0.1,
			},
		},
		{
			name: "Reset options",
			opts: []Options{
				WithConsumersOption(3),
				WithConsumersOption(2),
				WithConsumersOption(4),
				WithBufSize(3),
				WithBufSize(2),
				WithBufSize(4),
				WithBatchSize(3),
				WithBatchSize(2),
				WithBatchSize(4),
				WithDegradationFactor(0.2),
				WithDegradationFactor(0.3),
				WithDegradationFactor(0.1),
			},
			result: Options{
				bufSize:           4,
				consumers:         4,
				batchSize:         4,
				degradationFactor: 0.1,
			},
		},
		{
			name: "Empty options",
			opts: []Options{},
			result: Options{
				bufSize:           0,
				consumers:         1,
				batchSize:         1,
				degradationFactor: 0.75,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.EqualValues(t, tt.result, mergeOptions(tt.opts...))
		})
	}
}

func Test_MergeRangeReaderOptions(t *testing.T) {
	tests := []struct {
		name   string
		opts   []RangeReaderOptions
		result RangeReaderOptions
	}{
		{
			name: "Simple case",
			opts: []RangeReaderOptions{
				RangeReaderWithProducersOption(3),
				RangeReaderWithBatchSize(11),
			},
			result: RangeReaderOptions{
				producers: 3,
				batchSize: 11,
			},
		},
		{
			name: "Simple case reversed",
			opts: []RangeReaderOptions{
				RangeReaderWithProducersOption(3),
				RangeReaderWithBatchSize(11),
			},
			result: RangeReaderOptions{
				producers: 3,
				batchSize: 11,
			},
		},
		{
			name: "Reset options",
			opts: []RangeReaderOptions{
				RangeReaderWithProducersOption(3),
				RangeReaderWithProducersOption(2),
				RangeReaderWithProducersOption(4),
				RangeReaderWithBatchSize(3),
				RangeReaderWithBatchSize(2),
				RangeReaderWithBatchSize(4),
			},
			result: RangeReaderOptions{
				producers: 4,
				batchSize: 4,
			},
		},
		{
			name: "Empty options",
			opts: []RangeReaderOptions{},
			result: RangeReaderOptions{
				producers: 1,
				batchSize: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.EqualValues(t, tt.result, mergeRangeReaderOptions(tt.opts...))
		})
	}
}
