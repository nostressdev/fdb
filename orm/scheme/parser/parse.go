package parser

import (
	"strings"

	"github.com/nostressdev/fdb/errors"
	"github.com/nostressdev/fdb/orm/scheme"
)

type ValuesParser struct {
	Models map[string]scheme.Model
}

func (p *ValuesParser) parseField(value interface{}, fieldType string) interface{} {
	if value == nil {
		return nil
	}
	switch fieldType {
	case "int32":
		return int32(value.(int))
	case "int64":
		return int64(value.(int))
	case "uint32":
		return uint32(value.(int))
	case "uint64":
		return uint64(value.(int))
	case "string":
		return value.(string)
	case "bool":
		return value.(bool)
	case "float":
		return float32(value.(float64))
	case "double":
		return value.(float64)
	}
	if strings.HasPrefix(fieldType, "@") {
		typeName := fieldType[1:]
		if model, ok := p.Models[typeName]; ok {
			if newValue, ok := value.(map[string]interface{}); ok {
				return p.parseModel(newValue, model)
			}
			panic(errors.ParsingError.Newf("model %s: field %s is not a map", typeName, value))
		}
	}
	panic(errors.ParsingError.Newf("unknown type %s", fieldType))
}

func (p *ValuesParser) parseModel(value map[string]interface{}, model scheme.Model) map[string]interface{} {
	fieldsMap := value
	modelFieldNames := make(map[string]int)
	for i, field := range model.Fields {
		modelFieldNames[field.Name] = i
	}
	for name, value := range fieldsMap {
		if _, ok := modelFieldNames[name]; !ok {
			panic(errors.ParsingError.Newf("model %s: field %s is not defined", model.Name, name))
		}
		value := p.parseField(value, model.Fields[modelFieldNames[name]].Type)
		fieldsMap[name] = value
	}
	return fieldsMap
}

func (p *ValuesParser) parseModels(models []scheme.Model) {
	for _, model := range models {
		for _, field := range model.Fields {
			field.DefaultValue = p.parseField(field.DefaultValue, field.Type)
		}
	}
}

func (p *ValuesParser) parseTables(tables []scheme.Table) {
	for _, table := range tables {
		for _, column := range table.Columns {
			column.DefaultValue = p.parseField(column.DefaultValue, column.Type)
		}
	}
}

func (p *ValuesParser) ParseValues(config *scheme.GeneratorConfig) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = errors.InternalError.Newf("%v", r)
			}
		}
	}()
	p.parseModels(config.Models)
	p.parseTables(config.Tables)
	return err
}

func NewValuesParser(models []scheme.Model) (*ValuesParser, error) {
	parser := &ValuesParser{
		Models: make(map[string]scheme.Model),
	}
	for _, model := range models {
		if model.ExternalModel != "" {
			continue
		}
		if _, ok := parser.Models[model.Name]; ok {
			return nil, errors.ValidationError.Newf("model %s is duplicated", model.Name)
		}
		parser.Models[model.Name] = model
	}
	return parser, nil
}
