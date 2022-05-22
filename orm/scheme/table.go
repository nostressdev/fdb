package scheme

import "github.com/nostressdev/fdb/errors"

func (table *Table) validate() {
	if len(table.Columns) == 0 {
		panic(errors.ValidationError.Newf("table %s has no columns", table.Name))
	}
	if len(table.PK) == 0 {
		panic(errors.ValidationError.Newf("table %s has no primary key", table.Name))
	}
	columnsSet := table.validateColumns()
	table.validateIndexes(columnsSet)
	table.validatePK(columnsSet)
}

func (table *Table) validateColumns() map[string]bool {
	columnsSet := make(map[string]bool)
	for _, column := range table.Columns {
		if ok := columnsSet[column.Name]; ok {
			panic(errors.ValidationError.Newf("column %s:%s is duplicated", table.Name, column.Name))
		}
		columnsSet[column.Name] = true
		column.validate()
	}
	return columnsSet
}

func (table *Table) validateIndexes(columnsSet map[string]bool) {
	indexesSet := make(map[string]bool)
	for _, index := range table.RangeIndexes {
		if ok := indexesSet[index.Name]; ok {
			panic(errors.ValidationError.Newf("table %s: range index %s is duplicated", table.Name, index.Name))
		}
		index.validate(columnsSet)
	}
}

func (table *Table) validatePK(columnsSet map[string]bool) {
	pkSet := make(map[string]bool)
	for _, pk := range table.PK {
		if ok := columnsSet[pk]; !ok {
			panic(errors.ValidationError.Newf("table %s: primary key %s is not in columns", table.Name, pk))
		}
		if ok := pkSet[pk]; ok {
			panic(errors.ValidationError.Newf("table %s: primary key %s is duplicated", table.Name, pk))
		}
	}
}
