package scheme

import "github.com/nostressdev/fdb/errors"

func (column *Column) validate() {
	if column.Name == "" {
		panic(errors.ValidationError.Newf("table %s: column has no name", column.Table.Name))
	}
	if column.Type == "" {
		panic(errors.ValidationError.Newf("table %s: column %s has no type", column.Table.Name, column.Name))
	}
}
