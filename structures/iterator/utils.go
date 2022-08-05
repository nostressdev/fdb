package iterator

import "github.com/apple/foundationdb/bindings/go/src/fdb"

func getIterators(t fdb.ReadTransaction, reverse bool, rs ...fdb.Range) []Iterator {
	its := make([]Iterator, 0, len(rs))
	opts := fdb.RangeOptions{
		Reverse: reverse,
	}
	for _, r := range rs {
		its = append(its, iteratorFrom(t.GetRange(r, opts).Iterator()))
	}
	return its
}

func isKVEmpty(kv fdb.KeyValue) bool {
	return kv.Key == nil && kv.Value == nil
}
