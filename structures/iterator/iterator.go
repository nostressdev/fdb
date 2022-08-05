package iterator

import "github.com/apple/foundationdb/bindings/go/src/fdb"

type Iterator interface {
	Advance() bool
	Get() ([]byte, error) // get only values from fdb.KeyValue
}

type comparator func(a, b []byte) int // <0 means earlier, == 0 means equal

func MergeRanges(t fdb.ReadTransaction, reverse bool, c comparator, rs ...fdb.Range) (Iterator, error) {
	its := getIterators(t, reverse, rs...)
	return MergeIterators(c, its...)
}

func MergeIterators(c comparator, its ...Iterator) (Iterator, error) {
	return newMergeIterator(its, c)
}

func IntersectRanges(t fdb.ReadTransaction, reverse bool, c comparator, rs ...fdb.Range) (Iterator, error) {
	its := getIterators(t, reverse, rs...)
	return IntersectIterators(c, its...)
}

func IntersectIterators(c comparator, its ...Iterator) (Iterator, error) {
	return newIntersectIterator(its, c)
}