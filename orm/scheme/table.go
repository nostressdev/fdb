package scheme

import (
	"github.com/nostressdev/fdb/orm/scheme/utils"
)

func (table *Table) validate() {
	utils.Validatef(len(table.Columns) == 0, "table %s has no columns", table.Name)
	utils.Validatef(len(table.PK) == 0, "table %s has no primary key", table.Name)
	table.validateColumns()
	table.validateIndexes()
	table.validatePK()
}

func (table *Table) validateColumns() {
	columnsSet := make(map[string]bool)
	for _, column := range table.Columns {
		utils.Validatef(columnsSet[column.Name], "column %s:%s is duplicated", table.Name, column.Name)
		columnsSet[column.Name] = true
		column.validate()
	}
}

func (table *Table) validateIndexes() {
	indexesSet := make(map[string]bool)
	for _, index := range table.RangeIndexes {
		utils.Validatef(indexesSet[index.Name], "table %s: range index %s is duplicated", table.Name, index.Name)
		index.validate()
	}
}

func (table *Table) validatePK() {
	pkSet := make(map[string]bool)
	for _, pk := range table.PK {
		utils.Validatef(!table.ColumnsSet[pk], "table %s: primary key %s is not in columns", table.Name, pk)
		utils.Validatef(pkSet[pk], "table %s: primary key %s is duplicated", table.Name, pk)
	}
}
