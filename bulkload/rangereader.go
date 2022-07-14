package bulkload

import (
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"sync"
)

type RangeReader Reader[fdb.KeyValue]

func NewRangeReader(db fdb.Database, sub subspace.Subspace, opts ...RangeReaderOptions) RangeReader {
	options := mergeRangeReaderOptions(opts...)
	return func(ch chan fdb.KeyValue) error {
		begin, end := sub.FDBRangeKeys()
		keys, err := db.LocalityGetBoundaryKeys(fdb.KeyRange{Begin: begin, End: end}, options.producers-1, 0)
		if err != nil {
			return err
		}
		errCh := make(chan error)
		wg := new(sync.WaitGroup)
		keys = append(keys, end.FDBKey())
		wg.Add(len(keys))
		for len(keys) != 0 {
			next := keys[0]
			keys = keys[1:]
			nextSelector := fdb.KeySelector{Key: next}
			if len(keys) == 0 {
				nextSelector = fdb.FirstGreaterThan(next)
			}
			go func(begin fdb.KeySelector, end fdb.KeySelector) {
				tr, err := db.CreateTransaction()
				if err != nil {
					wg.Done()
					errCh <- err
					return
				}
				for {
					iter := tr.Snapshot().GetRange(fdb.SelectorRange{
						Begin: begin,
						End:   end,
					}, fdb.RangeOptions{Mode: fdb.StreamingModeWantAll, Limit: options.batchSize}).Iterator()
					cnt := 0
					lastPassed := fdb.Key{}
					for cnt < options.batchSize && iter.Advance() {
						cnt++
						if kv, err := iter.Get(); err != nil {
							wg.Done()
							errCh <- err
							return
						} else {
							ch <- kv
							lastPassed = kv.Key
						}
					}
					if cnt != options.batchSize {
						tr.Cancel()
						break
					}
					tr.Reset()
					begin = fdb.FirstGreaterThan(lastPassed)
				}
				wg.Done()
			}(fdb.KeySelector{Key: begin.FDBKey()}, nextSelector)
			begin = next
		}
		wg.Wait()
		close(errCh)
		return <-errCh
	}
}
