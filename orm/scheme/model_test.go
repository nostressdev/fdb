package scheme

import "testing"

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
			defer func() {
				if r := recover(); r != nil {
					if !tt.wantErr {
						t.Fatalf("panic: %v", r)
					}
					return
				}
				if tt.wantErr {
					t.Fatalf("want error, but no error")
				}
			}()
			model := &Model{
				Name:   tt.fields.Name,
				Fields: tt.fields.Fields,
			}
			model.validate()
		})
	}
}
