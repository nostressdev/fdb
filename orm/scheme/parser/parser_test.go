package parser

import (
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

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
			want: &scheme.GeneratorConfig{
				Models: []scheme.Model{
					{
						Name: "user",
						Fields: []scheme.Field{
							{Name: "id", Type: "string", DefaultValue: "field-default-id"},
							{Name: "name", Type: "string", DefaultValue: "field-default-name"},
						},
					},
					{
						Name: "profile",
						Fields: []scheme.Field{
							{Name: "user", Type: "@user", DefaultValue: struct {
								id   string
								name string
							}{id: "model-default-user-id", name: "model-default-user-name"}},
						},
					},
				},
				Tables: []scheme.Table{
					{
						Name:        "table",
						StoragePath: "table",
						RangeIndexes: []scheme.RangeIndex{
							{Name: "age", IK: []string{"age"}, Columns: []string{"age", "id"}, Async: true},
						},
						Columns: []scheme.Column{
							{Name: "id", Type: "string", DefaultValue: "column-default-id"},
							{Name: "age", Type: "uint32", DefaultValue: 20},
							{Name: "user", Type: "@user"},
						},
						PK: []string{"id"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text, err := ioutil.ReadFile(tt.filename)
			if err != nil {
				t.Fatalf("unable to read file %s: %v", tt.filename, err)
			}
			parser := NewParser()
			err = parser.AddYAMLReader(strings.NewReader(string(text)))
			if err != nil {
				t.Fatalf("unable to parse file %s: %v", tt.filename, err)
			}
			got, err := parser.GetConfig()
			if err != nil {
				t.Fatalf("GetConfig() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("GetConfig() = %v, want %v", got, tt.want)
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
			name:     "infinite loop test",
			filename: "../../../test/infinite-loop.yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text, err := ioutil.ReadFile(tt.filename)
			if err != nil {
				t.Fatalf("unable to read file %s: %v", tt.filename, err)
			}
			parser := NewParser()
			err = parser.AddYAMLReader(strings.NewReader(string(text)))
			if err != nil {
				t.Fatalf("unable to parse file %s: %v", tt.filename, err)
			}
			_, err = parser.GetConfig()
			if err == nil {
				t.Fatal("GetConfig() must return error")
			}
		})
	}
}
