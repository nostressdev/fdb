package scheme

import "testing"

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
		})
	}
}
