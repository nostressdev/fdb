package yaml

import (
	"github.com/nostressdev/fdb/orm/scheme"
	"github.com/nostressdev/fdb/orm/scheme/parser"
	"gopkg.in/yaml.v2"
)

func ParseYAML(text string) (*scheme.GeneratorConfig, error) {
	config := &scheme.GeneratorConfig{}
	err := yaml.Unmarshal([]byte(text), config)
	if err != nil {
		return nil, err
	}
	if err := config.Validate(); err != nil {
		return nil, err
	}
	valuesParser, err := parser.NewValuesParser(config.Models)
	if err != nil {
		return nil, err
	}
	valuesParser.ParseValues(config)
	return config, nil
}
