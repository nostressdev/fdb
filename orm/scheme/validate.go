package scheme

import "github.com/nostressdev/fdb/errors"

var primitives = map[string]struct{}{
	"int32":  {},
	"int64":  {},
	"uint32": {},
	"uint64": {},
	"string": {},
	"bool":   {},
	"float":  {},
	"double": {},
}

func (model *Model) validate() {
	if model.ExternalModel != "" {
		// TODO: validate external model
		return
	}
	set := make(map[string]struct{})
	for _, field := range model.Fields {
		if field.Type == "" {
			panic(errors.ValidationError.Newf("field %s:%s has no type", model.Name, field.Name))
		}
		if _, ok := set[field.Name]; ok {
			panic(errors.ValidationError.Newf("field %s:%s is duplicated", model.Name, field.Name))
		}
		set[field.Name] = struct{}{}
	}
}

func (c *GeneratorConfig) validateModels() {
	modelsSet := make(map[string]struct{})
	for _, model := range c.Models {
		if _, ok := modelsSet[model.Name]; ok {
			panic(errors.ValidationError.Newf("model %s is duplicated", model.Name))
		}
		modelsSet[model.Name] = struct{}{}
		model.validate()
	}
}

func (c *GeneratorConfig) validateColumns(table *Table) map[string]struct{} {
	columnsSet := make(map[string]struct{})
	for _, column := range table.Columns {
		if _, ok := columnsSet[column.Name]; ok {
			panic(errors.ValidationError.Newf("column %s:%s is duplicated", table.Name, column.Name))
		}
		columnsSet[column.Name] = struct{}{}
		if column.Name == "" {
			panic(errors.ValidationError.Newf("table %s: column has no name", table.Name))
		}
		if column.Type == "" {
			panic(errors.ValidationError.Newf("table %s: column %s has no type", table.Name, column.Name))
		}
	}
	for _, pk := range table.PK {
		if _, ok := columnsSet[pk]; !ok {
			panic(errors.ValidationError.Newf("table %s: primary key %s is not in columns", table.Name, pk))
		}
	}
	return columnsSet
}

func (c *GeneratorConfig) validateIndexes(table *Table, columnsSet map[string]struct{}) {
	indexesSet := make(map[string]struct{})
	for _, index := range table.RangeIndexes {
		if _, ok := indexesSet[index.Name]; ok {
			panic(errors.ValidationError.Newf("table %s: range index %s is duplicated", table.Name, index.Name))
		}
		if index.Name == "" {
			panic(errors.ValidationError.Newf("table %s: range index has no name", table.Name))
		}
		if len(index.IK) == 0 {
			panic(errors.ValidationError.Newf("table %s: range index %s has no ik", table.Name, index.Name))
		}
		if len(index.Columns) == 0 {
			panic(errors.ValidationError.Newf("table %s: range index %s has no columns", table.Name, index.Name))
		}
		for _, ik := range index.IK {
			if _, ok := columnsSet[ik]; !ok {
				panic(errors.ValidationError.Newf("table %s: range index %s: ik %s is not in columns", table.Name, index.Name, ik))
			}
		}
		for _, column := range index.Columns {
			if _, ok := columnsSet[column]; !ok {
				panic(errors.ValidationError.Newf("table %s: range index %s: column %s is not in columns", table.Name, index.Name, column))
			}
		}
	}
}

func (c *GeneratorConfig) validateTables() {
	tablesSet := make(map[string]struct{})
	for _, table := range c.Tables {
		if _, ok := tablesSet[table.Name]; ok {
			panic(errors.ValidationError.Newf("table %s is duplicated", table.Name))
		}
		tablesSet[table.Name] = struct{}{}
		if table.StoragePath == "" {
			panic(errors.ValidationError.Newf("table %s has no storage path", table.Name))
		}
		if len(table.Columns) == 0 {
			panic(errors.ValidationError.Newf("table %s has no columns", table.Name))
		}
		if len(table.PK) == 0 {
			panic(errors.ValidationError.Newf("table %s has no primary key", table.Name))
		}
		columnsSet := c.validateColumns(&table)
		c.validateIndexes(&table, columnsSet)
	}
}

func (c *GeneratorConfig) checkCycles() {
	names := make(map[string]*Model)
	for _, model := range c.Models {
		names[model.Name] = &model
	}
	used := make(map[*Model]int)
	stack := make([]*Model, 0)
	for _, model := range c.Models {
		if _, ok := used[&model]; !ok {
			stack = append(stack, &model)
			for len(stack) > 0 {
				model := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				used[model] = 1
				for _, field := range model.Fields {
					if _, ok := primitives[field.Type]; !ok {
						if nextModel, ok := names[field.Type[1:]]; !ok {
							panic(errors.ValidationError.Newf("model %s: field %s: type %s is not defined", nextModel.Name, field.Name, field.Type))
						} else if value := used[nextModel]; value == 1 {
							panic(errors.ValidationError.Newf("model %s: field %s: type %s is cyclic", model.Name, field.Name, field.Type))
						} else if value != 2 {
							stack = append(stack, nextModel)
						}
					}
				}
				used[model] = 2
			}
		}
	}
}

func (c *GeneratorConfig) Validate() (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = errors.InternalError.Newf("%v", r)
			}
		}
	}()
	c.validateModels()
	c.validateTables()
	c.checkCycles()
	return
}
