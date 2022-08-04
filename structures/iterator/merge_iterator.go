package iterator

type mergeIterator struct {
	its []Iterator
	vs  [][]byte
	end bool
	c   comparator
}

func newMergeIterator(its []Iterator, c comparator) (*mergeIterator, error) {
	vs := make([][]byte, 0)
	for _, it := range its {
		if it.Advance() {
			v, err := it.Get()
			if err != nil {
				return nil, err
			}
			vs = append(vs, v)
		} else {
			vs = append(vs, nil)
		}
	}

	it := &mergeIterator{
		its: its,
		vs:  vs,
		c:   c,
		end: true,
	}
	for _, v := range it.vs {
		if v != nil {
			it.end = false
			break
		}
	}

	return it, nil
}

func (it *mergeIterator) Advance() bool {
	return !it.end
}

func (it *mergeIterator) Get() ([]byte, error) {
	var res []byte
	k := -1

	for i, v := range it.vs {
		if v == nil {
			continue
		}
		if res == nil || it.c(v, res) < 0 {
			k = i
			res = v
		}
	}

	if it.its[k].Advance() {
		var err error
		it.vs[k], err = it.its[k].Get()
		if err != nil {
			return nil, err
		}
	} else {
		it.vs[k] = nil
	}

	it.end = true
	for _, v := range it.vs {
		if v != nil {
			it.end = false
			break
		}
	}
	return res, nil
}
