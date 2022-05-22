package scheme

import (
	"strings"

	"github.com/nostressdev/fdb/errors"
	"github.com/nostressdev/fdb/orm/scheme/utils"
	"github.com/nostressdev/fdb/orm/scheme/utils/graph"
)

var primitives = map[string]bool{
	"int32":  true,
	"int64":  true,
	"uint32": true,
	"uint64": true,
	"string": true,
	"bool":   true,
	"float":  true,
	"double": true,
}

func (c *GeneratorConfig) validateModels() {
	modelsSet := make(map[string]bool)
	for _, model := range c.Models {
		utils.Validatef(modelsSet[model.Name], "model %s is duplicated", model.Name)
		modelsSet[model.Name] = true
		model.validate()
	}
}

func (c *GeneratorConfig) validateTables() {
	tablesSet := make(map[string]bool)
	for _, table := range c.Tables {
		utils.Validatef(tablesSet[table.Name], "table %s is duplicated", table.Name)
		tablesSet[table.Name] = true
		table.validate()
	}
}

func (c *GeneratorConfig) checkCycles() {
	graph := graph.New()
	for _, model := range c.Models {
		graph.AddNode(model.Name)
		for _, field := range model.Fields {
			if primitives[field.Type] {
				continue
			}
			graph.AddEdge(model.Name, field.Type[1:])
		}
	}
	ok, cycle := graph.IsCyclic()
	utils.Validatef(ok, "models cycle detected: %s", strings.Join(cycle, " -> "))
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
