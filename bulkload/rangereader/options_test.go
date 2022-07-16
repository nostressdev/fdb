package rangereader

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MergeRangeReaderOptions(t *testing.T) {
	tests := []struct {
		name   string
		opts   []Options
		result Options
	}{
		{
			name: "Simple case",
			opts: []Options{
				WithProducersOption(3),
				WithBatchSize(11),
			},
			result: Options{
				producers: 3,
				batchSize: 11,
			},
		},
		{
			name: "Simple case reversed",
			opts: []Options{
				WithProducersOption(3),
				WithBatchSize(11),
			},
			result: Options{
				producers: 3,
				batchSize: 11,
			},
		},
		{
			name: "Reset options",
			opts: []Options{
				WithProducersOption(3),
				WithProducersOption(2),
				WithProducersOption(4),
				WithBatchSize(3),
				WithBatchSize(2),
				WithBatchSize(4),
			},
			result: Options{
				producers: 4,
				batchSize: 4,
			},
		},
		{
			name: "Empty options",
			opts: []Options{},
			result: Options{
				producers: 1,
				batchSize: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.EqualValues(t, tt.result, mergeOptions(tt.opts...))
		})
	}
}
