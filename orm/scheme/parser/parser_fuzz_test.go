package parser

import (
	"github.com/nostressdev/fdb/utils/errors"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func FuzzParseYaml(f *testing.F) {
	filenames := []string{
		"testdata/simple.yaml",
		"testdata/models-loop.yaml",
		"testdata/integer-limits.yaml",
	}
	for _, filename := range filenames {
		data, err := os.ReadFile(filename)
		if err != nil {
			f.Fatalf("failed to read testdata from (%s): %v", filename, err)
		}
		f.Add(string(data))
	}
	f.Fuzz(func(t *testing.T, orig string) {
		assert.NotPanics(t, func() {
			_, err := New().AddYAML(strings.NewReader(orig)).Parse()
			if err != nil {
				assert.Truef(t, errors.GetType(err) != errors.NoType, "unexpected error type: %v", err)
				assert.Truef(t, errors.GetType(err) != errors.InternalError, "unexpected error type: %v, %s", err)
			}
		})
	})
}

func FuzzParseJSON(f *testing.F) {
	filenames := []string{
		"testdata/integer-limits.json",
	}
	for _, filename := range filenames {
		data, err := os.ReadFile(filename)
		if err != nil {
			f.Fatalf("failed to read testdata from (%s): %v", filename, err)
		}
		f.Add(string(data))
	}
	f.Fuzz(func(t *testing.T, orig string) {
		assert.NotPanics(t, func() {
			_, err := New().AddJSON(strings.NewReader(orig)).Parse()
			if err != nil {
				assert.Truef(t, errors.GetType(err) != errors.NoType, "unexpected error type: %v", err)
				assert.Truef(t, errors.GetType(err) != errors.InternalError, "unexpected error type: %v", err)
			}
		})
	})
}
