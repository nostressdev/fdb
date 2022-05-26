package scheme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTable_validate(t *testing.T) {
	type fields struct {
		Name    string
		Columns []*Column
		PK      []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "empty name",
			fields: fields{
				Columns: []*Column{
					{
						Name: "A",
						Type: "int32",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "empty columns",
			fields: fields{
				Name: "name",
			},
			wantErr: true,
		},
		{
			name: "empty pk",
			fields: fields{
				Name: "name",
				Columns: []*Column{
					{
						Name: "A",
						Type: "int32",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "valid table",
			fields: fields{
				Name: "name",
				Columns: []*Column{
					{
						Name: "A",
						Type: "int32",
					},
					{
						Name: "B",
						Type: "string",
					},
					{
						Name: "C",
						Type: "float32",
					},
				},
				PK: []string{"B", "C"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() {
				table := &Table{
					Name:         tt.fields.Name,
					RangeIndexes: []*RangeIndex{},
					Columns:      tt.fields.Columns,
					PK:           tt.fields.PK,
					ColumnsSet:   make(map[string]bool),
				}
				for _, c := range table.Columns {
					table.ColumnsSet[c.Name] = true
					c.Table = table
				}
				table.validate()
			}
			if tt.wantErr {
				assert.Panics(t, f, "validate() should panic")
			} else {
				assert.NotPanics(t, f, "validate() should not panic")
			}
		})
	}
}
