package parser

import (
	"os"
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
			filename: "testdata/simple.yaml",
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
							{Name: "user", Type: "@user", DefaultValue: map[string]any{"id": string("field-default-id"), "name": string("field-default-name")}},
						},
						PK: []string{"id"},
					},
				},
			}),
		},
		{
			name:     "integer limits yaml test",
			filename: "testdata/integer-limits.yaml",
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
			filename: "testdata/integer-limits.json",
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
			name:     "default values in columns taken from model",
			filename: "testdata/default-values.yaml",
			want: FillValues(&scheme.GeneratorConfig{
				Models: []*scheme.Model{
					{
						Name: "user",
						Fields: []*scheme.Field{
							{Name: "name", Type: "string", DefaultValue: "Ivan"},
							{Name: "surname", Type: "string", DefaultValue: "Ivanov"},
						},
					},
				},
				Tables: []*scheme.Table{
					{
						Name: "users",
						Columns: []*scheme.Column{
							{Name: "user", Type: "@user", DefaultValue: map[string]interface{}{
								"name":    "Petya",
								"surname": "Ivanov",
							}},
						},
						PK: []string{"user"},
					},
				},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, err := os.Open(tt.filename)
			if err != nil {
				t.Fatalf("unable to read file %s: %v", tt.filename, err)
			}
			parser := New()
			parser.AddYAML(reader)
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
			filename: "testdata/models-loop.yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, err := os.Open(tt.filename)
			if err != nil {
				t.Fatalf("unable to read file %s: %v", tt.filename, err)
			}
			parser := New()
			parser.AddYAML(reader)
			_, err = parser.Parse()
			if err == nil && fdbErrors.GetType(err) == fdbErrors.ValidationError {
				t.Fatal("GetConfig() must return validation error")
			}
		})
	}
}
