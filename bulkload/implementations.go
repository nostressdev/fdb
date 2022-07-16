package bulkload

import (
	"golang.org/x/exp/constraints"
)

func RangeProducer[T constraints.Integer](begin T, end T) Producer[T] {
	return func(ch chan T) error {
		for begin < end {
			ch <- begin
			begin++
		}
		return nil
	}
}

func ListProducer[T constraints.Integer](len T) Producer[T] {
	return RangeProducer[T](0, len)
}

func ProducerFromList[T any](elements []T) Producer[T] {
	return func(ch chan T) error {
		for _, elem := range elements {
			ch <- elem
		}
		return nil
	}
}
