package iterator

import "github.com/apple/foundationdb/bindings/go/src/fdb"

type Iterator interface {
	Advance() bool
	Get() (fdb.KeyValue, error)
}

type Comparator func(kv1, kv2 fdb.KeyValue) int // <0 means earlier, == 0 means equal

func MergeRanges(t fdb.ReadTransaction, reverse bool, c Comparator, rs ...fdb.Range) (Iterator, error) {
	its := getIterators(t, reverse, rs...)
	return MergeIterators(c, its...)
}

func MergeIterators(c Comparator, its ...Iterator) (Iterator, error) {
	return newMergeIterator(its, c)
}

func IntersectRanges(t fdb.ReadTransaction, reverse bool, c Comparator, rs ...fdb.Range) (Iterator, error) {
	its := getIterators(t, reverse, rs...)
	return IntersectIterators(c, its...)
}

func IntersectIterators(c Comparator, its ...Iterator) (Iterator, error) {
	return newIntersectIterator(its, c)
}
