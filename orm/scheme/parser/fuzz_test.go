package parser

import (
	"os"
	"strings"
	"testing"

	"github.com/nostressdev/fdb/errors"
)

var filenames = []string{
	"test/simple.yaml",
	"test/models-loop.yaml",
	"test/integer-limits.yaml",
}

func FuzzParseYaml(f *testing.F) {
	println(os.Getwd())
	for _, filename := range filenames {
		data, err := os.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		f.Add(string(data))
	}
	f.Fuzz(func(t *testing.T, orig string) {
		parser := New()
		parser.AddYAML(strings.NewReader(orig))
		_, err := parser.Parse()
		if err != nil && errors.GetType(err) == errors.NoType {
			t.Errorf("ParseYAML(%q) = %v", orig, err)
		}
	})
}
