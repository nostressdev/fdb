package parser

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/nostressdev/fdb/errors"
	"github.com/nostressdev/fdb/orm/scheme"
	"github.com/nostressdev/fdb/orm/scheme/utils"
	"github.com/nostressdev/fdb/orm/scheme/utils/graph"
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
			if model == nil {
				panic(errors.ParsingError.Newf("model is not defined"))
			}
			if modelsSet[model.Name] {
				panic(errors.ParsingError.Newf("model %s: duplicated model name", model.Name))
			}
			p.Models[model.Name] = model
			modelsSet[model.Name] = true
		}
	}
	return configs
}

func (p *Parser) switchValue(value interface{}, fieldType string) interface{} {
	defer func() {
		if r := recover(); r != nil {
			panic(errors.ParsingError.Newf("cannot parse value %s to type %s", value, fieldType))
		}
	}()

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
	return nil
}

func (p *Parser) parseField(value interface{}, fieldType string) interface{} {
	if value == nil {
		return p.getDefaultValueFromType(fieldType)
	}
	res := p.switchValue(value, fieldType)
	if res != nil {
		return res
	}
	if len(fieldType) <= 1 {
		panic(errors.ParsingError.Newf("unknown type %s", fieldType))
	}
	if model, ok := p.Models[fieldType[1:]]; ok && strings.HasPrefix(fieldType, "@") {
		if structMap, ok := value.(map[string]interface{}); ok {
			return p.parseModelValues(structMap, model)
		} else if structInterfaceMap, ok := value.(map[interface{}]interface{}); ok {
			structMap := make(map[string]interface{})
			for k := range structInterfaceMap {
				if _, ok := k.(string); !ok {
					panic(errors.ParsingError.Newf("model %s: field %s is not defined", model.Name, k))
				}
				structMap[k.(string)] = structInterfaceMap[k]
			}
			return p.parseModelValues(structMap, model)
		}
		panic(errors.ParsingError.Newf("model %s: field %s is not a map", fieldType, value))
	}
	panic(errors.ParsingError.Newf("unknown type %s", fieldType))
}

func (p *Parser) getDefaultValueFromType(fieldType string) interface{} {
	switch fieldType {
	case "int32":
		return int32(0)
	case "int64":
		return int64(0)
	case "uint32":
		return uint32(0)
	case "uint64":
		return uint64(0)
	case "string":
		return ""
	case "bool":
		return false
	case "float":
		return float32(0)
	case "double":
		return float64(0)
	}
	if len(fieldType) <= 1 {
		panic(errors.ParsingError.Newf("unknown type %s", fieldType))
	}
	if model, ok := p.Models[fieldType[1:]]; ok && strings.HasPrefix(fieldType, "@") {
		fieldValues := make(map[string]interface{})
		for _, field := range model.Fields {
			if field == nil {
				panic(errors.ParsingError.Newf("model %s: field is not defined", model.Name))
			}
			fieldValues[field.Name] = field.DefaultValue
		}
		return fieldValues
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
		fieldsMap[name] = p.parseField(value, modelFieldNames[name].Type)
	}
	for _, field := range model.Fields {
		if fieldsMap[field.Name] == nil {
			fieldsMap[field.Name] = field.DefaultValue
		}
	}
	return fieldsMap
}

func (p *Parser) parseModel(model *scheme.Model) {
	for _, field := range model.Fields {
		if field == nil {
			panic(errors.ParsingError.Newf("model %s: field is not defined", model.Name))
		}
		field.DefaultValue = p.parseField(field.DefaultValue, field.Type)
	}
}

func (p *Parser) parseModels(models []*scheme.Model) {
	for _, model := range models {
		p.parseModel(model)
	}
}

func (p *Parser) parseTables(tables []*scheme.Table) {
	for _, table := range tables {
		if table == nil {
			panic(errors.ParsingError.Newf("table is not defined"))
		}
		for _, column := range table.Columns {
			if column == nil {
				panic(errors.ParsingError.Newf("table %s: column is not defined", table.Name))
			}
			column.DefaultValue = p.parseField(column.DefaultValue, column.Type)
		}
	}
}

func (p *Parser) parseValues(config *scheme.GeneratorConfig) {
	// p.parseModels(config.Models)
	p.parseTables(config.Tables)
}

func (p *Parser) validateModelTypes() {
	for _, model := range p.Models {
		for _, field := range model.Fields {
			if field == nil {
				panic(errors.ParsingError.Newf("model %s: field is not defined", model.Name))
			}
			if field.Type == "" {
				panic(errors.ParsingError.Newf("model %s: field %s has no type", model.Name, field.Name))
			}
			if !scheme.Primitives[field.Type] &&
				(!strings.HasPrefix(field.Type, "@") || len(field.Type) < 1 ||
					p.Models[field.Type[1:]] == nil) {
				panic(errors.ParsingError.Newf("model %s: field %s has unknown type %s", model.Name, field.Name, field.Type))
			}
		}
	}
}

func (p *Parser) getModelsOrder() []string {
	p.validateModelTypes()
	graph := graph.New()
	for _, model := range p.Models {
		graph.AddNode(model.Name)
		for _, field := range model.Fields {
			if scheme.Primitives[field.Type] {
				continue
			}
			graph.AddEdge(model.Name, field.Type[1:])
		}
	}
	order, ok := graph.TopSort()
	if !ok {
		hasCycle, cycle := graph.IsCyclic()
		if hasCycle {
			utils.Validatef(ok, "models cycle detected: %s", strings.Join(cycle, " -> "))
		} else {
			panic(errors.InternalError.New("unknown graph error"))
		}
	}
	return order
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
	order := p.getModelsOrder()
	for _, name := range order {
		model := p.Models[name]
		p.parseModel(model)
	}
	models := make([]*scheme.Model, 0)
	tables := make([]*scheme.Table, 0)
	for _, config := range configs {
		p.parseTables(config.Tables)
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
