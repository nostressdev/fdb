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

func (p *Parser) parseField(value interface{}, field scheme.Field) (interface{}, error) {
	switch field.Type {
	case "int32":
		return value.(int32), nil
	case "int64":
		return value.(int64), nil
	case "uint32":
		return value.(uint32), nil
	case "uint64":
		return value.(uint64), nil
	case "string":
		return value.(string), nil
	case "bool":
		return value.(bool), nil
	case "float":
		return value.(float32), nil
	case "double":
		return value.(float64), nil
	}
	if strings.HasPrefix(field.Type, "@") {
		typeName := field.Type[1:]
		if model, ok := p.Models[typeName]; ok {
			return p.parseModel(value, model)
		}
	}
	return nil, fmt.Errorf("unknown type %s", field.Type)
}

func (p *Parser) parseModel(value interface{}, model scheme.Model) (interface{}, error) {
	return nil, nil
}

func (p *Parser) parseConfig(text string) (*scheme.GeneratorConfig, error) {
	config := &scheme.GeneratorConfig{}
	err := yaml.Unmarshal([]byte(text), config)
	if err != nil {
		return nil, err
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
