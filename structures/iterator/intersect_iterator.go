package iterator

import "github.com/apple/foundationdb/bindings/go/src/fdb"

type intersectIterator struct { // necessary moveNext() when created
	its  []Iterator
	kvs  []fdb.KeyValue
	next fdb.KeyValue //synchronized with kvs
	c    Comparator
}

func newIntersectIterator(its []Iterator, c Comparator) (*intersectIterator, error) {
	kvs := make([]fdb.KeyValue, 0, len(its))
	for _, it := range its {
		if it.Advance() {
			kv, err := it.Get()
			if err != nil {
				return nil, err
			}
			kvs = append(kvs, kv)
		} else {
			return &intersectIterator{
				its:  its,
				next: fdb.KeyValue{},
				c:    c,
			}, nil
		}
	}
	it := &intersectIterator{
		its:  its,
		kvs:  kvs,
		next: fdb.KeyValue{},
		c:    c,
	}
	if err := it.moveNext(); err != nil {
		return nil, err
	}
	return it, nil
}

func (it *intersectIterator) Advance() bool {
	return !isKVEmpty(it.next)
}

func (it *intersectIterator) Get() (fdb.KeyValue, error) {
	res := it.next
	for i, iter := range it.its {
		if iter.Advance() {
			kv, err := iter.Get()
			if err != nil {
				return fdb.KeyValue{}, err
			}
			it.kvs[i] = kv
		} else {
			it.next = fdb.KeyValue{}
			return res, nil
		}
	}
	if err := it.moveNext(); err != nil {
		return fdb.KeyValue{}, err
	}
	return res, nil
}

func (it *intersectIterator) moveNext() error {
	var res fdb.KeyValue
	k := -1

	for i, kv := range it.kvs {
		if isKVEmpty(res) {
			res = kv
			k = i
			break
		}
	}
	if isKVEmpty(res) {
		it.next = fdb.KeyValue{}
		return nil
	}

	for i := 0; i < len(it.its); i++ {
		if i == k {
			continue
		}
		c := it.c(it.kvs[i], res)
		if c == 0 {
			continue
		}
		if c < 0 {
			for true {
				if it.its[i].Advance() {
					kv, err := it.its[i].Get()
					if err != nil {
						it.next = fdb.KeyValue{}
						return err
					}
					c = it.c(kv, res)
					if c >= 0 {
						it.kvs[i] = kv
						break
					}
				} else {
					it.next = fdb.KeyValue{}
					it.kvs[i] = fdb.KeyValue{}
					return nil
				}
			}
			if c == 0 {
				continue
			}
		}
		if c > 0 {
			res = it.kvs[i]
			k = i
			i = -1
		}
	}
	it.next = res
	return nil
}
