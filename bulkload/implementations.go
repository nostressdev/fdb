package bulkload

import "golang.org/x/exp/constraints"

func RangeGenerator[T constraints.Integer](begin T, end T) Reader[T] {
	return func(ch chan T) error {
		for begin < end {
			ch <- begin
			begin++
		}
		return nil
	}
}

func ListGenerator[T constraints.Integer](len T) Reader[T] {
	return RangeGenerator[T](0, len)
}
