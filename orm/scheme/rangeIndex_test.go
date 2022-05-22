package scheme

import "testing"

func TestRangeIndex_validate(t *testing.T) {
	table := &Table{
		Name: "table",
		ColumnsSet: map[string]bool{
			"A": true,
			"B": true,
			"C": true,
		},
	}
	type fields struct {
		Name    string
		IK      []string
		Columns []string
		Table   *Table
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "empty name",
			fields: fields{
				Table: table,
			},
			wantErr: true,
		},
		{
			name: "empty ik",
			fields: fields{
				Name:  "name",
				Table: table,
			},
			wantErr: true,
		},
		{
			name: "empty columns",
			fields: fields{
				Name:  "name",
				IK:    []string{"A", "B"},
				Table: table,
			},
			wantErr: true,
		},
		{
			name: "valid range index",
			fields: fields{
				Name:    "name",
				IK:      []string{"A", "B"},
				Columns: []string{"C"},
				Table:   table,
			},
			wantErr: false,
		},
		{
			name: "duplicated ik",
			fields: fields{
				Name:    "name",
				IK:      []string{"A", "B", "A"},
				Columns: []string{"C"},
				Table:   table,
			},
			wantErr: true,
		},
		{
			name: "duplicated columns",
			fields: fields{
				Name:    "name",
				IK:      []string{"A", "B"},
				Columns: []string{"C", "C"},
				Table:   table,
			},
			wantErr: true,
		},
		{
			name: "ik not in table",
			fields: fields{
				Name:    "name",
				IK:      []string{"A", "D", "C"},
				Columns: []string{"C"},
				Table:   table,
			},
			wantErr: true,
		},
		{
			name: "columns not in table",
			fields: fields{
				Name:    "name",
				IK:      []string{"A", "B"},
				Columns: []string{"C", "D"},
				Table:   table,
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
			index := &RangeIndex{
				Name:    tt.fields.Name,
				IK:      tt.fields.IK,
				Columns: tt.fields.Columns,
				Table:   tt.fields.Table,
			}
			index.validate()
		})
	}
}
