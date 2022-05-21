package parser

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/nostressdev/fdb/errors"
	"github.com/nostressdev/fdb/orm/scheme"
	"gopkg.in/yaml.v2"
)

type Parser struct {
	Models  map[string]scheme.Model
	Configs []*scheme.GeneratorConfig
}

type Decoder interface {
	Decode(v interface{}) error
}

func (p *Parser) AddDecoder(decoder Decoder) error {
	config := &scheme.GeneratorConfig{}
	err := decoder.Decode(config)
	if err != nil {
		return err
	}
	p.Configs = append(p.Configs, config)
	return nil
}

func (p *Parser) AddJSONReader(reader io.Reader) error {
	config := &scheme.GeneratorConfig{}
	err := json.NewDecoder(reader).Decode(config)
	if err != nil {
		return err
	}
	p.Configs = append(p.Configs, config)
	return nil
}

func (p *Parser) AddYAMLReader(reader io.Reader) error {
	config := &scheme.GeneratorConfig{}
	err := yaml.NewDecoder(reader).Decode(config)
	if err != nil {
		return err
	}
	p.Configs = append(p.Configs, config)
	return nil
}

func (p *Parser) init() {
	modelsSet := make(map[string]struct{})
	for _, config := range p.Configs {
		for _, model := range config.Models {
			if _, ok := modelsSet[model.Name]; ok {
				panic(errors.ParsingError.Newf("model %s: duplicated model name", model.Name))
			}
			p.Models[model.Name] = model
			modelsSet[model.Name] = struct{}{}
		}
	}
}

func (p *Parser) parseField(value interface{}, fieldType string) interface{} {
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

func (p *Parser) parseModel(value map[string]interface{}, model scheme.Model) map[string]interface{} {
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

func (p *Parser) parseModels(models []scheme.Model) {
	for _, model := range models {
		for _, field := range model.Fields {
			field.DefaultValue = p.parseField(field.DefaultValue, field.Type)
		}
	}
}

func (p *Parser) parseTables(tables []scheme.Table) {
	for _, table := range tables {
		for _, column := range table.Columns {
			column.DefaultValue = p.parseField(column.DefaultValue, column.Type)
		}
	}
}

func (p *Parser) parseValues(config *scheme.GeneratorConfig) (err error) {
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
	return
}

func (p *Parser) GetConfig() (*scheme.GeneratorConfig, error) {
	p.init()
	for _, config := range p.Configs {
		err := p.parseValues(config)
		if err != nil {
			return nil, err
		}
	}
	models := make([]scheme.Model, 0)
	tables := make([]scheme.Table, 0)
	for _, config := range p.Configs {
		models = append(models, config.Models...)
		tables = append(tables, config.Tables...)
	}
	return &scheme.GeneratorConfig{
		Models:  models,
		Tables:  tables,
	}, nil
}


func NewParser() (*Parser, error) {
	parser := &Parser{
		Models:  make(map[string]scheme.Model),
		Configs: []*scheme.GeneratorConfig{},
	}
	return parser, nil
}
