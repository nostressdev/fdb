package gen

import (
	"fmt"
	"github.com/nostressdev/fdb/orm/scheme"
	"strings"
)

func GenerateTable(gFile *GeneratedFile, table *scheme.Table, models []*scheme.Model) {
	generateStructs(gFile, table, models)
	generateMethods(gFile, table, models)
}

func generateStructs(gFile *GeneratedFile, table *scheme.Table, models []*scheme.Model) {
	generateTableStruct(gFile, table)
	generateTablePKStruct(gFile, table, models)
	generateTableRowStruct(gFile, table)
	generateFutureTableRowStruct(gFile, table)
}

func generateTableStruct(gFile *GeneratedFile, table *scheme.Table) {
	tableString :=
		`type %sTable struct {
			Enc lib.Encoder
			Dec lib.Decoder
			Sub subspace.Subspace
		}`
	newTableString :=
		`func New%[1]sTable(opts ...lib.TableOptions) (*%[1]sTable, error) {
			table := &%[1]sTable{}
			for _, opt := range opts {
				if opt.Enc != nil {
					table.Enc = opt.Enc
				}
				if opt.Dec != nil {
					table.Dec = opt.Dec
				}
				if opt.Sub != nil {
					table.Sub = opt.Sub
				}
			}
			if table.Enc == nil {
				return nil, fmt.Errorf("encoder is nil")
			}
			if table.Dec == nil {
				return nil, fmt.Errorf("decoder is nil")
			}
			if table.Sub == nil {
				return nil, fmt.Errorf("subspace is nil")
			}
			return table, nil
		}`

	gFile.Println(fmt.Sprintf(tableString, table.Name))
	gFile.Println(fmt.Sprintf(newTableString, table.Name))
}

func generateTablePKStruct(gFile *GeneratedFile, table *scheme.Table, models []*scheme.Model) {
	gFile.Println(fmt.Sprintf("type %sTablePK struct {", table.Name))
	for _, column := range table.PK {
		tp, err := getType(table, models, column)
		if err != nil {
			panic(err)
		}
		gFile.Println(fmt.Sprintf("	%s %s", strings.Join(strings.Split(column, "."), ""), tp))
	}
	gFile.Println("}")
	gFile.Println("")
}

func generateTableRowStruct(gFile *GeneratedFile, table *scheme.Table) {
	gFile.Println(fmt.Sprintf("type %sTableRow struct {", table.Name))
	for _, column := range table.Columns {
		gFile.Println(fmt.Sprintf("	%s %s", column.Name, column.Type))
	}
	gFile.Println("}")
	gFile.Println("")
}

func generateFutureTableRowStruct(gFile *GeneratedFile, table *scheme.Table) {
	futureTableRowString :=
		`type Future%sTableRow struct {
			Dec lib.Decoder
			Future fdb.FutureByteSlice
		}`
	gFile.Println(fmt.Sprintf(futureTableRowString, table.Name))
}

func generateMethods(gFile *GeneratedFile, table *scheme.Table, models []*scheme.Model) {
	generateTableMethods(gFile, table)
	generateTablePKMethods(gFile, table, models)
	generateFutureTableRowMethods(gFile, table)
}

func generateTableMethods(gFile *GeneratedFile, table *scheme.Table) {
	getString :=
		`func (table *%[1]sTable) Get(tr fdb.ReadTransaction, pk *%[1]sTablePK) (*Future%[1]sTableRow, error) {
			key, err := pk.Pack()
			if err != nil {
				return nil, err
			}
		
			future := tr.Get(table.Sub.Sub(key...))
			return &Future%[1]sTableRow{Future: future, Dec: table.Dec}, nil
		}`
	gFile.Println(fmt.Sprintf(getString, table.Name))

	mustGetString :=
		`func (table *%[1]sTable) MustGet(tr fdb.ReadTransaction, pk *%[1]sTablePK) *Future%[1]sTableRow {
			future, err := table.Get(tr, pk)
			if err != nil {
				panic(err)
			}
			return future
		}`
	gFile.Println(fmt.Sprintf(mustGetString, table.Name))

	insertString :=
		`func (table *%[1]sTable) Insert(tr fdb.Transaction, model *%[1]sTableRow) error {
			%[2]s
			key, err := pk.Pack()
			if err != nil {
				return err
			}
		
			value, err := table.Enc.Encode(model)
			if err != nil {
				return err
			}
			tr.Set(table.Sub.Sub(key...), value)
			return nil
		}`

	gFile.Println(fmt.Sprintf(insertString, table.Name, fillPK(table)))

	mustInsertString :=
		`func (table *%[1]sTable) MustInsert(tr fdb.Transaction, model *%[1]sTableRow) {
		err := table.Insert(tr, model)
		if err != nil {
			panic(err)
		}
	}`
	gFile.Println(fmt.Sprintf(mustInsertString, table.Name))

	deleteString :=
		`func (table *%[1]sTable) Delete(tr fdb.Transaction, pk *%[1]sTablePK) error {
			key, err := pk.Pack()
			if err != nil {
				return err
			}
			tr.Clear(table.Sub.Sub(key...))
			return nil
		}`
	gFile.Println(fmt.Sprintf(deleteString, table.Name))

	mustDeleteString :=
		`func (table *%[1]sTable) MustDelete(tr fdb.Transaction, pk *%[1]sTablePK) {
		err := table.Delete(tr, pk)
		if err != nil {
			panic(err)
		}
	}`
	gFile.Println(fmt.Sprintf(mustDeleteString, table.Name))
}

func fillPK(table *scheme.Table) string {
	var res string
	res += "pk := &" + table.Name + "TablePK{\n"
	for _, column := range table.PK {
		res += strings.Join(strings.Split(column, "."), "") + ": model." + column + ",\n"
	}
	res += "}"
	return res
}

func getType(table *scheme.Table, models []*scheme.Model, q string) (string, error) {
	if !strings.Contains(q, ".") {
		for _, column := range table.Columns {
			if column.Name == q {
				return column.Type, nil
			}
		}
		return "", fmt.Errorf("don't find \"%s\" in table columns", q)
	}

	qSlice := strings.Split(q, ".")        //делаем сплит
	for _, column := range table.Columns { // ищем колонку по названию
		if column.Name == qSlice[0] {
			for _, model := range models { // ищем модель, которая имеет название как наш тип
				if column.Type == model.Name {
					for _, field := range model.Fields { // ищем поле
						if field.Name == qSlice[1] {
							return field.Type, nil
						}
					}
					return "", fmt.Errorf("don't find field \"%s\" from \"%s\" in model \"%s\"", qSlice[1], q, model.Name)
				}
			}
			return "", fmt.Errorf("don't find model \"%s\" as type of column \"%s\" from \"%s\"", column.Type, column.Name, q)
		}
	}
	return "", fmt.Errorf("don't find \"%s\" from \"%s\" in table columns", qSlice[0], q)
}

func generateTablePKMethods(gFile *GeneratedFile, table *scheme.Table, models []*scheme.Model) {
	gFile.Println("func (pk *" + table.Name + "TablePK) Pack() ([]tuple.TupleElement, error) {")
	gFile.Println("	var err error")
	elements := make([]string, 0, len(table.PK))
	for _, column := range table.PK {
		name := strings.Join(strings.Split(column, "."), "")
		tp, err := getType(table, models, column)
		if err != nil {
			panic(err)
		}
		switch tp {
		case "string":
			stringPackSting := `pk%[1]sBytes := []byte(pk.%[1]s)`
			gFile.Println(fmt.Sprintf(stringPackSting, name))
		case "uint64":
			uint64PackString :=
				`%[1]sBuf := new(bytes.Buffer)
				err = binary.Write(%[1]sBuf, binary.BigEndian, pk.%[1]s)
				if err != nil {
					return nil, err
				}
				pk%[1]sBytes := %[1]sBuf.Bytes()`
			gFile.Println(fmt.Sprintf(uint64PackString, name))
		case "int64":
			int64PackString :=
				`%[1]sBuf := new(bytes.Buffer)
				err = binary.Write(%[1]sBuf, binary.BigEndian, pk.%[1]s)
				if err != nil {
					return nil, err
				}
				pk%[1]sBytes := %[1]sBuf.Bytes()`
			gFile.Println(fmt.Sprintf(int64PackString, name))
		case "float32":
			float32PackString :=
				`%[1]sBuf := new(bytes.Buffer)
				err = binary.Write(%[1]sBuf, binary.BigEndian, pk.%[1]s)
				if err != nil {
					return nil, err
				}
				pk%[1]sBytes := %[1]sBuf.Bytes()`
			gFile.Println(fmt.Sprintf(float32PackString, name))
		default:
			panic("unknown type " + tp)
		}
		elements = append(elements, "pk"+name+"Bytes")
	}
	gFile.Println("	if err != nil {")
	gFile.Println("		return nil, err")
	gFile.Println("	}")
	gFile.Println("	return []tuple.TupleElement{" + strings.Join(elements, ", ") + "}, nil")
	gFile.Println("}")
}

func generateFutureTableRowMethods(gFile *GeneratedFile, table *scheme.Table) {
	funcNewTableRowString :=
		`func (future *Future%[1]sTableRow) new%[1]sTableRow(value []byte) (*%[1]sTableRow, error) {
			if row, err := future.Dec.Decode(value); err != nil {
				return nil, err
			} else {
				return row.(*%[1]sTableRow), nil
			}
		}`
	gFile.Println(fmt.Sprintf(funcNewTableRowString, table.Name))

	funcGetString :=
		`func (future *Future%[1]sTableRow) Get() (*%[1]sTableRow, error) {
		value, err := future.Future.Get()
		if err != nil {
			return nil, err
		}
		return future.new%[1]sTableRow(value)
	}`
	gFile.Println(fmt.Sprintf(funcGetString, table.Name))

	funcMustGetString :=
		`func (future *Future%[1]sTableRow) MustGet() *%[1]sTableRow {
		value, err := future.Get()
		if err != nil {
			panic(err)
		}
		return value
	}`
	gFile.Println(fmt.Sprintf(funcMustGetString, table.Name))
}
