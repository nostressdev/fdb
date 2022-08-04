package iterator

import "github.com/apple/foundationdb/bindings/go/src/fdb"

type fdbIterator struct {
	it *fdb.RangeIterator
}

func iteratorFrom(i *fdb.RangeIterator) *fdbIterator {
	return &fdbIterator{it: i}
}

func (i *fdbIterator) Advance() bool {
	return i.it.Advance()
}

func (i *fdbIterator) Get() ([]byte, error) {
	kv, err := i.it.Get()
	if err != nil {
		return nil, err
	}
	return kv.Value, err
}
