package testsGenerate

import (
	"github.com/nostressdev/fdb/gen"
	"github.com/nostressdev/fdb/tests"
	"testing"
)

func Test_Generate(t *testing.T) {
	gen.GenFiles(tests.Config)
}
