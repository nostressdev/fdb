package rangereader

type Options struct {
	batchSize int
	producers int
}

func WithProducersOption(value int) Options {
	return Options{
		producers: value,
	}
}

func WithBatchSize(value int) Options {
	return Options{
		batchSize: value,
	}
}

func mergeOptions(options ...Options) Options {
	result := Options{
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
