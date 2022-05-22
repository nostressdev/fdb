package scheme

import "github.com/nostressdev/fdb/orm/scheme/utils"

func (column *Column) validate() {
	utils.Validatef(column.Name == "", "table %s: column has no name", column.Table.Name)
	utils.Validatef(column.Type == "", "table %s: column %s has no type", column.Table.Name, column.Name)
}
