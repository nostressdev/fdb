package tests

import (
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/nostressdev/fdb/lib"
	gen_fdb "github.com/nostressdev/fdb/tests/generated"
)

var toAge = lib.TableOptions{Enc: &gen_fdb.AgeSortTableRowJsonEncoder{}, Dec: &gen_fdb.AgeSortTableRowJsonDecoder{}, Sub: subspace.Sub("tests")}
var toUsers = lib.TableOptions{Enc: &gen_fdb.UsersTableRowJsonEncoder{}, Dec: &gen_fdb.UsersTableRowJsonDecoder{}, Sub: subspace.Sub("tests")}
