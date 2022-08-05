package iterator

type intersectIterator struct { // necessary moveNext() when created
	its  []Iterator
	vs   [][]byte
	next []byte //synchronized with vs
	c    comparator
}

func newIntersectIterator(its []Iterator, c comparator) (*intersectIterator, error) {
	vs := make([][]byte, 0, len(its))
	for _, it := range its {
		if it.Advance() {
			v, err := it.Get()
			if err != nil {
				return nil, err
			}
			vs = append(vs, v)
		} else {
			return &intersectIterator{
				its:  its,
				next: nil,
				c:    c,
			}, nil
		}
	}
	it := &intersectIterator{
		its:  its,
		vs:   vs,
		next: nil,
		c:    c,
	}
	if err := it.moveNext(); err != nil {
		return nil, err
	}
	return it, nil
}

func (it *intersectIterator) Advance() bool {
	return it.next != nil
}

func (it *intersectIterator) Get() ([]byte, error) {
	res := it.next
	for i, iter := range it.its {
		if iter.Advance() {
			v, err := iter.Get()
			if err != nil {
				return nil, err
			}
			it.vs[i] = v
		} else {
			it.next = nil
			return res, nil
		}
	}
	if err := it.moveNext(); err != nil {
		return nil, err
	}
	return res, nil
}

func (it *intersectIterator) moveNext() error {
	var res []byte
	k := -1

	for i, v := range it.vs {
		if res == nil {
			res = v
			k = i
			break
		}
	}
	if res == nil {
		it.next = nil
		return nil
	}

	for i := 0; i < len(it.its); i++ {
		if i == k {
			continue
		}
		c := it.c(it.vs[i], res)
		if c == 0 {
			continue
		}
		if c < 0 {
			for true {
				if it.its[i].Advance() {
					v, err := it.its[i].Get()
					if err != nil {
						it.next = nil
						return err
					}
					c = it.c(v, res)
					if c >= 0 {
						it.vs[i] = v
						break
					}
				} else {
					it.next = nil
					it.vs[i] = nil
					return nil
				}
			}
			if c == 0 {
				continue
			}
		}
		if c > 0 {
			res = it.vs[i]
			k = i
			i = -1
		}
	}
	it.next = res
	return nil
}
