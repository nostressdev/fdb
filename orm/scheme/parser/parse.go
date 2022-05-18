package parser

import (
	"fmt"
	"strings"

	"github.com/nostressdev/fdb/orm/scheme"
	"github.com/nostressdev/fdb/orm/scheme/errors"
)

type ValuesParser struct {
	Models map[string]scheme.Model
}

func (p *ValuesParser) parseField(value interface{}, fieldType string) (interface{}, error) {
	switch fieldType {
	case "int32":
		return int32(value.(int)), nil
	case "int64":
		return int64(value.(int)), nil
	case "uint32":
		return uint32(value.(int)), nil
	case "uint64":
		return uint64(value.(int)), nil
	case "string":
		return value.(string), nil
	case "bool":
		return value.(bool), nil
	case "float":
		return float32(value.(float64)), nil
	case "double":
		return value.(float64), nil
	}
	if strings.HasPrefix(fieldType, "@") {
		typeName := fieldType[1:]
		if model, ok := p.Models[typeName]; ok {
			return p.parseModel(value, model)
		}
	}
	return nil, fmt.Errorf("unknown type %s", fieldType)
}

func (p *ValuesParser) parseModel(value interface{}, model scheme.Model) (interface{}, error) {
	var fieldsMap map[string]interface{}
	var ok bool
	if fieldsMap, ok = value.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("model %s must be a map", model.Name)
	}
	modelFieldNames := make(map[string]int)
	for i, field := range model.Fields {
		modelFieldNames[field.Name] = i
	}
	for name, value := range fieldsMap {
		if _, ok := modelFieldNames[name]; !ok {
			return nil, fmt.Errorf("model %s: field %s is not defined", model.Name, name)
		}
		value, err := p.parseField(value, model.Fields[modelFieldNames[name]].Type)
		if err != nil {
			return nil, err
		}
		fieldsMap[name] = value
	}
	return fieldsMap, nil
}

func (p *ValuesParser) parseModels(models []scheme.Model) error {
	var err error
	for _, model := range models {
		for _, field := range model.Fields {
			if field.DefaultValue != nil {
				if field.DefaultValue, err = p.parseField(field.DefaultValue, field.Type); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (p *ValuesParser) parseTables(tables []scheme.Table) error {
	var err error
	for _, table := range tables {
		for _, column := range table.Columns {
			if column.DefaultValue != nil {
				if column.DefaultValue, err = p.parseField(column.DefaultValue, column.Type); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (p *ValuesParser) ParseValues(config *scheme.GeneratorConfig) error {
	if err := p.parseModels(config.Models); err != nil {
		return err
	}
	return p.parseTables(config.Tables)
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
