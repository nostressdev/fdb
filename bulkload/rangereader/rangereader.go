package rangereader

import (
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/nostressdev/fdb/bulkload"
	"github.com/nostressdev/fdb/utils/future"
	"golang.org/x/xerrors"
)

type RangeReader bulkload.Producer[fdb.KeyValue]

func readSubRange(tr fdb.Transaction, tasks chan fdb.KeyValue, rangeStart fdb.KeySelector, rangeEnd fdb.KeySelector, options Options) error {
	for {
		cnt := 0
		iter := tr.Snapshot().GetRange(fdb.SelectorRange{
			Begin: rangeStart,
			End:   rangeEnd,
		}, fdb.RangeOptions{
			Mode:  fdb.StreamingModeWantAll,
			Limit: options.batchSize,
		}).Iterator()
		hasElements := false
		for iter.Advance() {
			hasElements = true
			if kv, err := iter.Get(); err != nil {
				var trErr fdb.Error
				if xerrors.As(err, &trErr) {
					if !((trErr.Code >= 2101 && trErr.Code <= 2103) || trErr.Code == 1007) {
						if err = tr.OnError(trErr).Get(); err != nil {
							return err
						}
					}
					tr.Reset()
					break
				}
				return err
			} else {
				cnt++
				rangeStart = fdb.FirstGreaterThan(kv.Key)
				tasks <- kv
				if cnt == options.batchSize {
					tr.Reset()
					break
				}
			}
		}
		if !hasElements {
			return nil
		}
	}
}

func processIterFunction(db fdb.Database, tasks chan fdb.KeyValue, rangeStart fdb.Key, rangeEnd fdb.Key, options Options) func() error {
	return func() error {
		tr, err := db.CreateTransaction()
		if err != nil {
			return err
		}
		return readSubRange(tr, tasks, fdb.FirstGreaterOrEqual(rangeStart), fdb.FirstGreaterOrEqual(rangeEnd), options)
	}
}

func NewRangeReader(db fdb.Database, kr fdb.KeyRange, opts ...Options) RangeReader {
	options := mergeOptions(opts...)
	return func(tasks chan fdb.KeyValue) error {
		var keys []fdb.Key
		if options.producers > 1 {
			var err error
			keys, err = db.LocalityGetBoundaryKeys(kr, options.producers-1, 0)
			if err != nil {
				return err
			}
		}
		keys = append(keys, kr.End.FDBKey())
		futures := make([]future.FutureNil, 0, len(keys))
		rangeStart := kr.Begin.FDBKey()
		for len(keys) != 0 {
			rangeEnd := keys[0]
			keys = keys[1:]
			futures = append(futures, future.AsyncNil(processIterFunction(
				db, tasks, rangeStart, rangeEnd, options,
			)))
			rangeStart = rangeEnd
		}
		return future.AwaitAllNil(futures...)
	}
}
