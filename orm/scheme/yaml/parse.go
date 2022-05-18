package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/nostressdev/fdb/orm/scheme"
	"gopkg.in/yaml.v3"
)

type Parser struct {
	Models map[string]scheme.Model
}

func (p *Parser) parseField(value interface{}, fieldType string) (interface{}, error) {

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

func (p *Parser) parseModel(value interface{}, model scheme.Model) (interface{}, error) {
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

func (p *Parser) initModels(models []scheme.Model) error {
	p.Models = make(map[string]scheme.Model)
	for _, model := range models {
		if model.ExternalModel != "" {
			continue
		}
		if _, ok := p.Models[model.Name]; ok {
			return fmt.Errorf("model %s is duplicated", model.Name)
		}
		p.Models[model.Name] = model
	}
	return nil
}

func (p *Parser) parseConfig(text string) (*scheme.GeneratorConfig, error) {
	config := &scheme.GeneratorConfig{}
	err := yaml.Unmarshal([]byte(text), config)
	if err != nil {
		return nil, err
	}
	if err := config.Validate(); err != nil {
		return nil, err
	}
	p.initModels(config.Models)
	for _, model := range config.Models {
		for _, field := range model.Fields {
			if field.DefaultValue != nil {
				field.DefaultValue, err = p.parseField(field.DefaultValue, field.Type)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	for _, table := range config.Tables {
		for _, column := range table.Columns {
			if column.DefaultValue != nil {
				column.DefaultValue, err = p.parseField(column.DefaultValue, column.Type)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return config, nil
}

func main() {
	if len(os.Args) != 2 {
		panic("Usage: go run parse.go config.yaml")
	}
	// read text from file
	text, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	p := &Parser{}
	config, err := p.parseConfig(string(text))
	if err != nil {
		panic(err)
	}
	fmt.Println(config)
}
