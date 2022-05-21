package test

import (
	"strings"
	"testing"

	"github.com/nostressdev/fdb/errors"
	"github.com/nostressdev/fdb/orm/scheme/parser"
)

func FuzzParseYaml(f *testing.F) {
	f.Fuzz(func(t *testing.T, orig string) {
		parser := parser.NewParser()
		err := parser.AddYAMLReader(strings.NewReader(orig))
		if err != nil && errors.GetType(err) != errors.NoType {
			t.Errorf("ParseYAML(%q) = %v", orig, err)
		}
		_, err = parser.GetConfig()
		if err != nil && errors.GetType(err) != errors.NoType {
			t.Errorf("ParseYAML(%q) = %v", orig, err)
		}
	})
}