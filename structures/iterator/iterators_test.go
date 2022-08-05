package iterator

import (
	"bytes"
	"fmt"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/stretchr/testify/require"
	"log"
	"math/rand"
	"sort"
	"testing"
	"time"
)

const (
	cntIterators = 100
	cntValues    = 100
	chance       = 0.98
)

func Test_Iterators(t *testing.T) {
	rand.Seed(time.Now().Unix())
	require.NoError(t, fdb.APIVersion(600), "use fdb 6.3.22")
	db, err := fdb.OpenDefault()
	require.NoError(t, err, "open foundationdb")
	sub := subspace.Sub("Test_Iterators")

	cleanup := func(t *testing.T) {
		_, err := db.Transact(func(tr fdb.Transaction) (interface{}, error) {
			tr.ClearRange(sub)
			return nil, nil
		})
		require.NoError(t, err)
	}

	cmp := func(kv1, kv2 fdb.KeyValue) int {
		return bytes.Compare(kv1.Value, kv2.Value)
	}

	t.Run("cleanup start", cleanup)

	ss := make([]subspace.Subspace, 0, cntIterators)
	rs := make([]fdb.Range, 0, cntIterators)
	for i := 0; i < cntIterators; i++ {
		ss = append(ss, sub.Sub(i))
		rs = append(rs, ss[len(ss)-1])
	}
	res := make(map[byte]int)
	t.Run("init", func(t *testing.T) {
		_, err := db.Transact(func(tr fdb.Transaction) (interface{}, error) {
			for _, s := range ss {
				for i := 0; i < cntValues; i++ {
					if rand.Float64() < chance {
						tr.Set(s.Sub([]byte{byte('a' + i)}), []byte{byte('a' + i)})
						res[byte('a'+i)]++
					}
				}
			}
			return nil, nil
		})
		require.NoError(t, err)
	})

	log.Println(res)
	keys := make([]byte, 0, len(res))
	for k := range res {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	t.Run("merge", func(t *testing.T) {
		_, err := db.ReadTransact(func(tr fdb.ReadTransaction) (interface{}, error) {
			it, err := MergeRanges(tr, false, cmp, rs...)
			require.NoError(t, err, "create iterator")
			for _, k := range keys {
				log.Println("k =", k)
				if it.Advance() {
					kv, err := it.Get()
					require.NoError(t, err, "get iterator")
					log.Println("kv =", kv)
					require.Equal(t, k, kv.Value[0], "equal")
				} else {
					require.NoError(t, fmt.Errorf("iterator not advance"))
				}
			}
			if it.Advance() {
				kv, err := it.Get()
				require.NoError(t, err, "get iterator last")
				require.NoError(t, fmt.Errorf("more get in iterator, value=%s", string(kv.Value)))
			}
			return nil, nil
		})
		require.NoError(t, err)
	})

	keys = make([]byte, 0)
	for k, v := range res {
		if v != cntIterators {
			continue
		}
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	t.Run("intersect", func(t *testing.T) {
		_, err := db.ReadTransact(func(tr fdb.ReadTransaction) (interface{}, error) {
			it, err := IntersectRanges(tr, false, cmp, rs...)
			require.NoError(t, err, "create iterator")
			for _, k := range keys {
				log.Println("k =", k)
				if it.Advance() {
					kv, err := it.Get()
					require.NoError(t, err, "get iterator")
					log.Println("kv =", kv)
					require.Equal(t, k, kv.Value[0], "equal")
				} else {
					require.NoError(t, fmt.Errorf("iterator not advance"))
				}
			}
			if it.Advance() {
				kv, err := it.Get()
				require.NoError(t, err, "get iterator last")
				require.NoError(t, fmt.Errorf("more get in iterator, value=%s", string(kv.Value)))
			}
			return nil, nil
		})
		require.NoError(t, err)
	})
	t.Run("intersect+merge", func(t *testing.T) {
		_, err := db.ReadTransact(func(tr fdb.ReadTransaction) (interface{}, error) {
			it1, err := IntersectRanges(tr, false, cmp, rs...)
			require.NoError(t, err, "create iterator")
			it2, err := MergeRanges(tr, false, cmp, rs...)
			require.NoError(t, err, "create iterator")
			it, err := IntersectIterators(cmp, it1, it2)
			require.NoError(t, err, "create iterator")
			for _, k := range keys {
				log.Println("k =", k)
				if it.Advance() {
					kv, err := it.Get()
					require.NoError(t, err, "get iterator")
					log.Println("kv =", kv)
					require.Equal(t, k, kv.Value[0], "equal")
				} else {
					require.NoError(t, fmt.Errorf("iterator not advance"))
				}
			}
			if it.Advance() {
				kv, err := it.Get()
				require.NoError(t, err, "get iterator last")
				require.NoError(t, fmt.Errorf("more get in iterator, value=%s", string(kv.Value)))
			}
			return nil, nil
		})
		require.NoError(t, err)
	})

	t.Run("intersect+merge with empty", func(t *testing.T) {
		_, err := db.ReadTransact(func(tr fdb.ReadTransaction) (interface{}, error) {
			var r fdb.Range
			r = sub.Sub("hello")
			it1, err := IntersectRanges(tr, false, cmp, append(rs, r)...)
			require.NoError(t, err, "create iterator")
			it2, err := MergeRanges(tr, false, cmp, append(rs, r)...)
			require.NoError(t, err, "create iterator")
			it, err := IntersectIterators(cmp, it1, it2)
			require.NoError(t, err, "create iterator")
			if it.Advance() {
				kv, err := it.Get()
				require.NoError(t, err, "get iterator last")
				require.NoError(t, fmt.Errorf("more get in iterator, value=%s", string(kv.Value)))
			}
			return nil, nil
		})
		require.NoError(t, err)
	})

	t.Run("intersect 0", func(t *testing.T) {
		_, err := db.ReadTransact(func(tr fdb.ReadTransaction) (interface{}, error) {
			it, err := IntersectRanges(tr, false, cmp)
			require.NoError(t, err, "create iterator")
			if it.Advance() {
				kv, err := it.Get()
				require.NoError(t, err, "get iterator last")
				require.NoError(t, fmt.Errorf("more get in iterator, value=%s", string(kv.Value)))
			}
			return nil, nil
		})
		require.NoError(t, err)
	})

	t.Run("intersect 1", func(t *testing.T) {
		_, err := db.Transact(func(tr fdb.Transaction) (interface{}, error) {
			tr.Set(sub.Sub("hello").Sub("hello"), []byte("hello"))
			it, err := IntersectRanges(tr, false, cmp, sub.Sub("hello"))
			require.NoError(t, err, "create iterator")
			if it.Advance() {
				kv, err := it.Get()
				require.NoError(t, err, "get iterator")
				log.Println("kv =", kv)
				require.Equal(t, "hello", string(kv.Value), "equal")
			} else {
				require.NoError(t, fmt.Errorf("iterator not advance"))
			}
			if it.Advance() {
				kv, err := it.Get()
				require.NoError(t, err, "get iterator last")
				require.NoError(t, fmt.Errorf("more get in iterator, value=%s", string(kv.Value)))
			}
			return nil, nil
		})
		require.NoError(t, err)
	})

	t.Run("cleanup end", cleanup)
}
