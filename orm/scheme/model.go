package scheme

import "github.com/nostressdev/fdb/orm/scheme/utils"

func (model *Model) validate() {
	if model.ExternalModel != "" {
		// TODO: validate external model
		return
	}
	set := make(map[string]bool)
	for _, field := range model.Fields {
		utils.Validatef(field.Type == "", "field %s:%s has no type", model.Name, field.Name)
		utils.Validatef(set[field.Name], "field %s:%s is duplicated", model.Name, field.Name)
		set[field.Name] = true
	}
}
