package scheme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColumn_validate(t *testing.T) {
	type fields struct {
		Name         string
		Type         string
		DefaultValue interface{}
		Table        *Table
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "empty name & type",
			fields: fields{
				Table: &Table{},
			},
			wantErr: true,
		},
		{
			name: "empty type",
			fields: fields{
				Name:  "name",
				Table: &Table{},
			},
			wantErr: true,
		},
		{
			name: "valid column",
			fields: fields{
				Name:  "name",
				Type:  "type",
				Table: &Table{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() {
				column := &Column{
					Name:         tt.fields.Name,
					Type:         tt.fields.Type,
					DefaultValue: tt.fields.DefaultValue,
					Table:        tt.fields.Table,
				}
				column.validate()
			}
			if tt.wantErr {
				assert.Panics(t, f, "validate() should panic")
			} else {
				assert.NotPanics(t, f, "validate() should not panic")
			}
		})
	}
}
