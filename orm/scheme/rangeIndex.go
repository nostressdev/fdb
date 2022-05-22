package scheme

import "github.com/nostressdev/fdb/orm/scheme/utils"

func (index *RangeIndex) validate() {
	utils.Validatef(index.Name == "", "table %s: range index has no name", index.Table.Name)
	utils.Validatef(len(index.IK) == 0, "table %s: range index %s has no ik", index.Table.Name, index.Name)
	utils.Validatef(len(index.Columns) == 0, "table %s: range index %s has no columns", index.Table.Name, index.Name)
	ikSet := make(map[string]bool)
	for _, ik := range index.IK {
		utils.Validatef(!index.Table.ColumnsSet[ik], "table %s: range index %s: ik %s is not in columns", index.Table.Name, index.Name, ik)
		utils.Validatef(ikSet[ik], "table %s: range index %s: ik %s is duplicated", index.Table.Name, index.Name, ik)

	}
	indexColumnsSet := make(map[string]bool)
	for _, column := range index.Columns {
		utils.Validatef(!index.Table.ColumnsSet[column], "table %s: range index %s: column %s is not in columns", index.Table.Name, index.Name, column)
		utils.Validatef(indexColumnsSet[column], "table %s: range index %s: column %s is duplicated", index.Table.Name, index.Name, column)
	}
}
