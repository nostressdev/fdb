package scheme

import "testing"

func TestGeneratorConfig_Validate(t *testing.T) {
	type fields struct {
		Models []*Model
		Tables []*Table
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "duplicating models",
			fields: fields{
				Models: []*Model{
					{
						Name: "A",
					},
					{
						Name: "B",
					},
					{
						Name: "A",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "duplicating tables",
			fields: fields{
				Tables: []*Table{
					{
						Name: "A",
					},
					{
						Name: "B",
					},
					{
						Name: "A",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "cycled models",
			fields: fields{
				Models: []*Model{
					{
						Name: "A",
						Fields: []*Field{
							{
								Name: "B",
								Type: "@B",
							},
						},
					},
					{
						Name: "B",
						Fields: []*Field{
							{
								Name: "A",
								Type: "@A",
							},
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &GeneratorConfig{
				Models: tt.fields.Models,
				Tables: tt.fields.Tables,
			}
			if err := c.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("GeneratorConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
