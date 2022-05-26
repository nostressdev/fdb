package scheme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModel_validate(t *testing.T) {
	type fields struct {
		Name   string
		Fields []*Field
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "empty name",
			fields: fields{
				Fields: []*Field{},
			},
			wantErr: true,
		},
		{
			name: "empty fields",
			fields: fields{
				Name: "name",
			},
			wantErr: false,
		},
		{
			name: "valid model",
			fields: fields{
				Name: "name",
				Fields: []*Field{
					{
						Name: "name",
						Type: "type",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "duplicated field",
			fields: fields{
				Name: "name",
				Fields: []*Field{
					{
						Name: "A",
						Type: "type",
					},
					{
						Name: "B",
						Type: "int32",
					},
					{
						Name: "A",
						Type: "string",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() {
				model := &Model{
					Name:   tt.fields.Name,
					Fields: tt.fields.Fields,
				}
				model.validate()
			}
			if tt.wantErr {
				assert.Panics(t, f, "validate() should panic")
			} else {
				assert.NotPanics(t, f, "validate() should not panic")
			}
		})
	}
}
