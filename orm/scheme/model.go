package scheme

import "github.com/nostressdev/fdb/errors"

func (model *Model) validate() {
	if model.ExternalModel != "" {
		// TODO: validate external model
		return
	}
	set := make(map[string]bool)
	for _, field := range model.Fields {
		if field.Type == "" {
			panic(errors.ValidationError.Newf("field %s:%s has no type", model.Name, field.Name))
		}
		if ok := set[field.Name]; ok {
			panic(errors.ValidationError.Newf("field %s:%s is duplicated", model.Name, field.Name))
		}
		set[field.Name] = true
	}
}
