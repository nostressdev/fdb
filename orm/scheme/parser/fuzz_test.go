package parser

import (
	"strings"
	"testing"

	"github.com/nostressdev/fdb/errors"
)

func FuzzParseYaml(f *testing.F) {
	f.Fuzz(func(t *testing.T, orig string) {
		parser := New()
		err := parser.AddYAML(strings.NewReader(orig))
		if err != nil && errors.GetType(err) != errors.NoType {
			t.Errorf("ParseYAML(%q) = %v", orig, err)
		}
		_, err = parser.Parse()
		if err != nil && errors.GetType(err) != errors.NoType {
			t.Errorf("ParseYAML(%q) = %v", orig, err)
		}
	})
}
