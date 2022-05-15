package fdb

type GeneratorConfig struct {
	Models []Model
	Tables []Table
}

type Model struct {
	Name          string
	Fields        map[string]Field
	ExternalModel string // продумать как хранить инфу о протобафе, может что то дополнительное
}

type Field struct {
	Name         string
	Type         string
	DefaultValue interface{}
}

type Table struct {
	Name        string
	StoragePath string
	Indexes     []Index
	Columns     map[string]Column
	PK          []string // column names
}

type Column struct {
	Name         string
	Type         string
	DefaultValue interface{}
}

type Index struct {
	Name    string
	IK      []string // column names
	Columns []string
	Async   bool
}
