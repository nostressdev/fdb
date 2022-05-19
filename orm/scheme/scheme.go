package scheme

type GeneratorConfig struct {
	Models []*Model
	Tables []*Table
}

type Model struct {
	Name          string
	Fields        []Field
	ExternalModel string // продумать как хранить инфу о протобафе, может что то дополнительное
}

type Field struct {
	Name         string
	Type         string
	DefaultValue interface{}
}

type Table struct {
	Name         string
	StoragePath  string
	RangeIndexes []RangeIndex
	Columns      []Column
	PK           []string // column names
}

type Column struct {
	Name         string
	Type         string
	DefaultValue interface{}
}

type RangeIndex struct {
	Name    string
	IK      []string // column names
	Columns []string
	Async   bool
}
