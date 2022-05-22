package scheme

import "github.com/nostressdev/fdb/errors"

func (index *RangeIndex) validate(columnsSet map[string]bool) {
	if index.Name == "" {
		panic(errors.ValidationError.Newf("table %s: range index has no name", index.Table.Name))
	}
	if len(index.IK) == 0 {
		panic(errors.ValidationError.Newf("table %s: range index %s has no ik", index.Table.Name, index.Name))
	}
	if len(index.Columns) == 0 {
		panic(errors.ValidationError.Newf("table %s: range index %s has no columns", index.Table.Name, index.Name))
	}
	ikSet := make(map[string]bool)
	for _, ik := range index.IK {
		if ok := columnsSet[ik]; !ok {
			panic(errors.ValidationError.Newf("table %s: range index %s: ik %s is not in columns", index.Table.Name, index.Name, ik))
		}
		if ok := ikSet[ik]; ok {
			panic(errors.ValidationError.Newf("table %s: range index %s: ik %s is duplicated", index.Table.Name, index.Name, ik))
		}

	}
	indexColumnsSet := make(map[string]bool)
	for _, column := range index.Columns {
		if ok := columnsSet[column]; !ok {
			panic(errors.ValidationError.Newf("table %s: range index %s: column %s is not in columns", index.Table.Name, index.Name, column))
		}
		if ok := indexColumnsSet[column]; ok {
			panic(errors.ValidationError.Newf("table %s: range index %s: column %s is duplicated", index.Table.Name, index.Name, column))
		}
	}
}
