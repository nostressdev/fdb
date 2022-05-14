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
	Name string
	Type string // продумать как его задавать, может быть enum
}

type Table struct {
	Name       string
	Columns    map[string]Column
	PrimaryKey []string // column names
}

type Column struct {
	Name         string
	Type         string // продумать как его задавать, может быть enum
	DefaultValue interface{}
}
