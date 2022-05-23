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
	Models   map[string]*scheme.Model
	Decoders []decoder
}

type decoder interface {
	Decode(v interface{}) error
}

func (p *Parser) addDecoder(d decoder) {
	p.Decoders = append(p.Decoders, d)
}

func (p *Parser) AddJSON(reader io.Reader) *Parser {
	p.addDecoder(json.NewDecoder(reader))
	return p
}

func (p *Parser) AddYAML(reader io.Reader) *Parser {
	p.addDecoder(yaml.NewDecoder(reader))
	return p
}

func (p *Parser) init() []*scheme.GeneratorConfig {
	modelsSet := make(map[string]bool)
	configs := make([]*scheme.GeneratorConfig, 0)
	for _, decoder := range p.Decoders {
		config := &scheme.GeneratorConfig{}
		err := decoder.Decode(config)
		if err != nil {
			panic(errors.ParsingError.Wrap(err, "parsing config"))
		}
		configs = append(configs, config)
		for i := range config.Models {
			model := config.Models[i]
			if modelsSet[model.Name] {
				panic(errors.ParsingError.Newf("model %s: duplicated model name", model.Name))
			}
			p.Models[model.Name] = model
			modelsSet[model.Name] = true
		}
	}
	return configs
}

func (p *Parser) parseField(value interface{}, fieldType string) interface{} {
	if value == nil {
		return nil
	}
	// While parsing default values we don't have exact type of field,
	// only `interface{}`, so we need to check it here.
	// gopkg.in/yaml.v2 and encoding/json Unmarshall interpret all integer values
	// (including signed and unsigned) as `int`, and if value is out of bounds
	// of `int`, it will be interpreted as `int64` and if it is out of bounds of
	// `int64` it will be interpreted as `uint64`.
	// Also `int` bounds depend on architecture, so we need to handle both cases.
	switch fieldType {
	case "int32":
		return int32(value.(int))
	case "int64":
		if val, ok := value.(int); ok {
			return int64(val)
		} else if val, ok := value.(int64); ok {
			return val
		}
	case "uint32":
		if val, ok := value.(int); ok {
			return uint32(val)
		} else if val, ok := value.(int64); ok {
			return uint32(val)
		}
	case "uint64":
		return value.(uint64)
	case "string":
		return value.(string)
	case "bool":
		return value.(bool)
	case "float":
		return float32(value.(float64))
	case "double":
		return value.(float64)
	}
	if model, ok := p.Models[fieldType[1:]]; ok && strings.HasPrefix(fieldType, "@") {
		if structMap, ok := value.(map[string]interface{}); ok {
			return p.parseModelValues(structMap, model)
		} else if structInterfaceMap, ok := value.(map[interface{}]interface{}); ok {
			structMap := make(map[string]interface{})
			for k := range structInterfaceMap {
				structMap[k.(string)] = structInterfaceMap[k]
			}
			return p.parseModelValues(structMap, model)
		}
		panic(errors.ParsingError.Newf("model %s: field %s is not a map", fieldType, value))
	}
	panic(errors.ParsingError.Newf("unknown type %s", fieldType))
}

func (p *Parser) parseModelValues(fieldsMap map[string]interface{}, model *scheme.Model) map[string]interface{} {
	modelFieldNames := make(map[string]*scheme.Field)
	for _, field := range model.Fields {
		modelFieldNames[field.Name] = field
	}
	for name, value := range fieldsMap {
		if _, ok := modelFieldNames[name]; !ok {
			panic(errors.ParsingError.Newf("model %s: field %s is not defined", model.Name, name))
		}
		value := p.parseField(value, modelFieldNames[name].Type)
		fieldsMap[name] = value
	}
	return fieldsMap
}

func (p *Parser) parseModels(models []*scheme.Model) {
	for _, model := range models {
		for _, field := range model.Fields {
			field.DefaultValue = p.parseField(field.DefaultValue, field.Type)
		}
	}
}

func (p *Parser) parseTables(tables []*scheme.Table) {
	for _, table := range tables {
		for _, column := range table.Columns {
			column.DefaultValue = p.parseField(column.DefaultValue, column.Type)
		}
	}
}

func (p *Parser) parseValues(config *scheme.GeneratorConfig) {
	p.parseModels(config.Models)
	p.parseTables(config.Tables)
}

func (p *Parser) Parse() (config *scheme.GeneratorConfig, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				if errors.GetType(e) == errors.NoType {
					err = errors.InternalError.New(e.Error())
				} else {
					err = e
				}
			} else {
				err = errors.InternalError.Newf("%v", r)
			}
		}
	}()
	configs := p.init()
	models := make([]*scheme.Model, 0)
	tables := make([]*scheme.Table, 0)
	for _, config := range configs {
		p.parseValues(config)
		models = append(models, config.Models...)
		tables = append(tables, config.Tables...)
	}
	config = &scheme.GeneratorConfig{
		Models: models,
		Tables: tables,
	}
	FillValues(config)
	return config, config.Validate()
}

func FillValues(config *scheme.GeneratorConfig) *scheme.GeneratorConfig {
	for _, table := range config.Tables {
		table.ColumnsSet = make(map[string]bool)
		for _, column := range table.Columns {
			column.Table = table
			table.ColumnsSet[column.Name] = true
		}
		for _, index := range table.RangeIndexes {
			index.Table = table
		}
	}
	return config
}

func New() *Parser {
	parser := &Parser{
		Models: make(map[string]*scheme.Model),
	}
	return parser
}
