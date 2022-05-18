package scheme

import "fmt"

type GeneratorConfig struct {
	Models []Model `yaml:"models"`
	Tables []Table `yaml:"tables"`
}

type Model struct {
	Name          string  `yaml:"name"`
	Fields        []Field `yaml:"fields"`
	ExternalModel string  `yaml:"external-model"`
}

type Field struct {
	Name         string      `yaml:"name"`
	Type         string      `yaml:"type"`
	DefaultValue interface{} `yaml:"default-value"`
}

type Table struct {
	Name         string       `yaml:"name"`
	StoragePath  string       `yaml:"storage-path"`
	RangeIndexes []RangeIndex `yaml:"range-indexes"`
	Columns      []Column     `yaml:"columns"`
	PK           []string     `yaml:"pk"`
}

type Column struct {
	Name         string      `yaml:"name"`
	Type         string      `yaml:"type"`
	DefaultValue interface{} `yaml:"default-value"`
}

type RangeIndex struct {
	Name    string   `yaml:"name"`
	IK      []string `yaml:"ik"`
	Columns []string `yaml:"columns"`
	Async   bool     `yaml:"async"`
}

func (c *GeneratorConfig) Validate() error {
	modelsSet := make(map[string]struct{})
	for _, model := range c.Models {
		if _, ok := modelsSet[model.Name]; ok {
			return fmt.Errorf("model %s is duplicated", model.Name)
		}
		modelsSet[model.Name] = struct{}{}
		if model.ExternalModel != "" {
			continue
		}
		set := make(map[string]struct{})
		for _, field := range model.Fields {
			if field.Type == "" {
				return fmt.Errorf("field %s:%s has no type", model.Name, field.Name)
			}
			if _, ok := set[field.Name]; ok {
				return fmt.Errorf("field %s:%s is duplicated", model.Name, field.Name)
			}
			set[field.Name] = struct{}{}
		}

	}
	// validate tables
	tablesSet := make(map[string]struct{})
	for _, table := range c.Tables {
		if _, ok := tablesSet[table.Name]; ok {
			return fmt.Errorf("table %s is duplicated", table.Name)
		}
		tablesSet[table.Name] = struct{}{}
		if table.StoragePath == "" {
			return fmt.Errorf("table %s has no storage path", table.Name)
		}
		if len(table.Columns) == 0 {
			return fmt.Errorf("table %s has no columns", table.Name)
		}
		if len(table.PK) == 0 {
			return fmt.Errorf("table %s has no primary key", table.Name)
		}
		columnsSet := make(map[string]struct{})
		for _, column := range table.Columns {
			if _, ok := columnsSet[column.Name]; ok {
				return fmt.Errorf("column %s:%s is duplicated", table.Name, column.Name)
			}
			columnsSet[column.Name] = struct{}{}
			if column.Name == "" {
				return fmt.Errorf("table %s: column has no name", table.Name)
			}
			if column.Type == "" {
				return fmt.Errorf("table %s: column %s has no type", table.Name, column.Name)
			}
		}
		for _, pk := range table.PK {
			if _, ok := columnsSet[pk]; !ok {
				return fmt.Errorf("table %s: primary key %s is not in columns", table.Name, pk)
			}
		}
		indexesSet := make(map[string]struct{})
		for _, index := range table.RangeIndexes {
			if _, ok := indexesSet[index.Name]; ok {
				return fmt.Errorf("table %s: range index %s is duplicated", table.Name, index.Name)
			}
			if index.Name == "" {
				return fmt.Errorf("table %s: range index has no name", table.Name)
			}
			if len(index.IK) == 0 {
				return fmt.Errorf("table %s: range index %s has no ik", table.Name, index.Name)
			}
			if len(index.Columns) == 0 {
				return fmt.Errorf("table %s: range index %s has no columns", table.Name, index.Name)
			}
			for _, ik := range index.IK {
				if _, ok := columnsSet[ik]; !ok {
					return fmt.Errorf("table %s: range index %s: ik %s is not in columns", table.Name, index.Name, ik)
				}
			}
			for _, column := range index.Columns {
				if _, ok := columnsSet[column]; !ok {
					return fmt.Errorf("table %s: range index %s: column %s is not in columns", table.Name, index.Name, column)
				}
			}
		}
	}

	return nil
}
