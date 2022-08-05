package iterator

import "github.com/apple/foundationdb/bindings/go/src/fdb"

type mergeIterator struct {
	its []Iterator
	kvs []fdb.KeyValue
	end bool
	c   Comparator
}

func newMergeIterator(its []Iterator, c Comparator) (*mergeIterator, error) {
	kvs := make([]fdb.KeyValue, 0, len(its))
	for _, it := range its {
		if it.Advance() {
			kv, err := it.Get()
			if err != nil {
				return nil, err
			}
			kvs = append(kvs, kv)
		} else {
			kvs = append(kvs, fdb.KeyValue{})
		}
	}

	it := &mergeIterator{
		its: its,
		kvs: kvs,
		c:   c,
		end: true,
	}
	for _, kv := range it.kvs {
		if !isKVEmpty(kv) {
			it.end = false
			break
		}
	}

	return it, nil
}

func (it *mergeIterator) Advance() bool {
	return !it.end
}

func (it *mergeIterator) Get() (fdb.KeyValue, error) {
	var res fdb.KeyValue

	for _, kv := range it.kvs {
		if isKVEmpty(kv) {
			continue
		}
		if isKVEmpty(res) || it.c(kv, res) < 0 {
			res = kv
		}
	}

	for i, iter := range it.its {
		if isKVEmpty(it.kvs[i]) || it.c(it.kvs[i], res) != 0 {
			continue
		}
		if iter.Advance() {
			var err error
			it.kvs[i], err = iter.Get()
			if err != nil {
				return fdb.KeyValue{}, err
			}
		} else {
			it.kvs[i] = fdb.KeyValue{}
		}
	}

	it.end = true
	for _, kv := range it.kvs {
		if !isKVEmpty(kv) {
			it.end = false
			break
		}
	}
	return res, nil
}
