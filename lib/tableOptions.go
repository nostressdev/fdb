package lib

import "github.com/apple/foundationdb/bindings/go/src/fdb/subspace"

type TableOptions struct {
	Enc Encoder
	Sub subspace.Subspace
}
