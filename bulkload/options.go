package bulkload

type Options struct {
	consumers int
	bufSize   int
	batchSize int
}

func WithConsumersOption(value int) Options {
	return Options{
		consumers: value,
	}
}

func WithBufSize(value int) Options {
	return Options{
		bufSize: value,
	}
}
func WithBatchSize(value int) Options {
	return Options{
		batchSize: value,
	}
}

func mergeOptions(options ...Options) Options {
	result := Options{
		consumers: 1,
		bufSize:   0,
		batchSize: 1,
	}
	for _, option := range options {
		if result.consumers < option.consumers {
			result.consumers = option.consumers
		}
		if result.bufSize < option.bufSize {
			result.bufSize = option.bufSize
		}
		if result.batchSize < option.batchSize {
			result.batchSize = option.batchSize
		}
	}
	return result
}

type RangeReaderOptions struct {
	batchSize int
	producers int
}

func RangeReaderWithProducersOption(value int) RangeReaderOptions {
	return RangeReaderOptions{
		producers: value,
	}
}

func RangeReaderWithBatchSize(value int) RangeReaderOptions {
	return RangeReaderOptions{
		batchSize: value,
	}
}

func mergeRangeReaderOptions(options ...RangeReaderOptions) RangeReaderOptions {
	result := RangeReaderOptions{
		producers: 1,
		batchSize: 1,
	}
	for _, option := range options {
		if result.producers < option.producers {
			result.producers = option.producers
		}
		if result.batchSize < option.batchSize {
			result.batchSize = option.batchSize
		}
	}
	return result
}
