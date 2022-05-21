package tests

import (
	"github.com/nostressdev/fdb/gen"
	"testing"
)

func Test_Generate(t *testing.T) {
	gen.GenFiles(Config)
}
