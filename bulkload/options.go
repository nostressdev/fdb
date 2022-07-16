package bulkload

type Options struct {
	consumers         int
	bufSize           int
	batchSize         int
	degradationFactor float64
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

func WithDegradationFactor(value float64) Options {
	return Options{
		degradationFactor: value,
	}
}

func mergeOptions(options ...Options) Options {
	result := Options{
		consumers:         1,
		bufSize:           0,
		batchSize:         1,
		degradationFactor: 0.75,
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
		if option.degradationFactor != 0 {
			result.degradationFactor = option.degradationFactor
		}
	}
	return result
}
