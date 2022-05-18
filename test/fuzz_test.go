package test

import (
	"testing"

	"github.com/nostressdev/fdb/orm/scheme/errors"
	"github.com/nostressdev/fdb/orm/scheme/yaml"
)

func FuzzParseYaml(f *testing.F) {
	f.Fuzz(func(t *testing.T, orig string) {
		_, err := yaml.ParseYAML(orig)
		if err != nil && errors.GetType(err) != errors.NoType {
			t.Errorf("ParseYAML(%q) = %v", orig, err)
		}
	})
}
