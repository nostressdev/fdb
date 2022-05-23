package tests

import (
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	generated "github.com/nostressdev/fdb/gen/tests/generated"
	"github.com/nostressdev/fdb/lib"
)

var toAge = lib.TableOptions{Enc: &generated.AgeSortTableRowJsonEncoder{}, Dec: &generated.AgeSortTableRowJsonDecoder{}, Sub: subspace.Sub("tests")}
var toUsers = lib.TableOptions{Enc: &generated.UsersTableRowJsonEncoder{}, Dec: &generated.UsersTableRowJsonDecoder{}, Sub: subspace.Sub("tests")}
