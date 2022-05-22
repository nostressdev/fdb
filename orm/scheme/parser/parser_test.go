package parser

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	fdbErrors "github.com/nostressdev/fdb/errors"
	"github.com/nostressdev/fdb/orm/scheme"
)

func TestParseYAML(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     *scheme.GeneratorConfig
	}{
		{
			name:     "simple test",
			filename: "../../../test/simple.yaml",
			want: FillValues(&scheme.GeneratorConfig{
				Models: []*scheme.Model{
					{
						Name: "profile",
						Fields: []*scheme.Field{
							{Name: "user", Type: "@user", DefaultValue: map[string]interface{}{"id": "model-default-user-id", "name": "model-default-user-name"}},
						},
					},
					{
						Name: "user",
						Fields: []*scheme.Field{
							{Name: "id", Type: "string", DefaultValue: "field-default-id"},
							{Name: "name", Type: "string", DefaultValue: "field-default-name"},
						},
					},
					{
						Name:          "external",
						ExternalModel: "filename.proto/MessageName",
					},
				},
				Tables: []*scheme.Table{
					{
						Name: "table",
						RangeIndexes: []*scheme.RangeIndex{
							{Name: "age", IK: []string{"age"}, Columns: []string{"age", "id"}, Async: true},
						},
						Columns: []*scheme.Column{
							{Name: "id", Type: "string", DefaultValue: "column-default-id"},
							{Name: "age", Type: "uint32", DefaultValue: uint32(20)},
							{Name: "user", Type: "@user"},
						},
						PK: []string{"id"},
					},
				},
			}),
		},
		{
			name:     "integer limits yaml test",
			filename: "../../../test/integer-limits.yaml",
			want: FillValues(&scheme.GeneratorConfig{
				Models: []*scheme.Model{
					{
						Name: "a",
						Fields: []*scheme.Field{
							{Name: "int64", Type: "int64", DefaultValue: int64(9223372036854775807)},
							{Name: "uint64", Type: "uint64", DefaultValue: uint64(18446744073709551615)},
							{Name: "int32", Type: "int32", DefaultValue: int32(2147483647)},
							{Name: "uint32", Type: "uint32", DefaultValue: uint32(4294967295)},
						},
					},
				},
				Tables: []*scheme.Table{},
			}),
		},
		{
			name:     "integer limits json test",
			filename: "../../../test/integer-limits.json",
			want: FillValues(&scheme.GeneratorConfig{
				Models: []*scheme.Model{
					{
						Name: "a",
						Fields: []*scheme.Field{
							{Name: "int64", Type: "int64", DefaultValue: int64(9223372036854775807)},
							{Name: "uint64", Type: "uint64", DefaultValue: uint64(18446744073709551615)},
							{Name: "int32", Type: "int32", DefaultValue: int32(2147483647)},
							{Name: "uint32", Type: "uint32", DefaultValue: uint32(4294967295)},
						},
					},
				},
				Tables: []*scheme.Table{},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text, err := ioutil.ReadFile(tt.filename)
			if err != nil {
				t.Fatalf("unable to read file %s: %v", tt.filename, err)
			}
			parser := New()
			parser.AddYAML(strings.NewReader(string(text)))
			got, err := parser.Parse()
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Fatalf("Parse() diff = %v", diff)
			}
		})
	}
}

func TestParseYAMLWithErrors(t *testing.T) {
	tests := []struct {
		name     string
		filename string
	}{
		{
			name:     "models loop test",
			filename: "../../../test/models-loop.yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text, err := ioutil.ReadFile(tt.filename)
			if err != nil {
				t.Fatalf("unable to read file %s: %v", tt.filename, err)
			}
			parser := New()
			parser.AddYAML(strings.NewReader(string(text)))
			_, err = parser.Parse()
			if err == nil && fdbErrors.GetType(err) == fdbErrors.ValidationError {
				t.Fatal("GetConfig() must return validation error")
			}
			fmt.Println(err.Error())
		})
	}
}
