package scheme

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
