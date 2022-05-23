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
	gFile.Println("type " + table.Name + "Table struct {")
	gFile.Println("	Enc lib.Encoder")
	gFile.Println("	Dec lib.Decoder")
	gFile.Println("	Sub subspace.Subspace")
	gFile.Println("}")
	gFile.Println("")
	gFile.Println("func New" + table.Name + "Table(opts ...lib.TableOptions) (*" + table.Name + "Table, error) {")
	gFile.Println("	table := &" + table.Name + "Table{}")
	gFile.Println("	for _, opt := range opts {")
	gFile.Println("		if opt.Enc != nil {")
	gFile.Println("			table.Enc = opt.Enc")
	gFile.Println("		}")
	gFile.Println("		if opt.Dec != nil {")
	gFile.Println("			table.Dec = opt.Dec")
	gFile.Println("		}")
	gFile.Println("		if opt.Sub != nil {")
	gFile.Println("			table.Sub = opt.Sub")
	gFile.Println("		}")
	gFile.Println("	}")
	gFile.Println("	if table.Enc == nil {")
	gFile.Println("		return nil, fmt.Errorf(\"encoder is nil\")")
	gFile.Println("	}")
	gFile.Println("	if table.Dec == nil {")
	gFile.Println("		return nil, fmt.Errorf(\"decoder is nil\")")
	gFile.Println("	}")
	gFile.Println("	if table.Sub == nil {")
	gFile.Println("		return nil, fmt.Errorf(\"subspace is nil\")")
	gFile.Println("	}")
	gFile.Println("	return table, nil")
	gFile.Println("}")
	gFile.Println("")
}

func generateTablePKStruct(gFile *GeneratedFile, table *scheme.Table, models []*scheme.Model) {
	gFile.Println("type " + table.Name + "TablePK struct {")
	for _, column := range table.PK {
		tp, err := getType(table, models, column)
		if err != nil {
			panic(err)
		}
		gFile.Println("		" + strings.Join(strings.Split(column, "."), "") + " " + tp)
	}
	gFile.Println("	}")
	gFile.Println("")
}

func generateTableRowStruct(gFile *GeneratedFile, table *scheme.Table) {
	gFile.Println("type " + table.Name + "TableRow struct {")
	for _, column := range table.Columns {
		gFile.Println("		" + column.Name + " " + column.Type)
	}
	gFile.Println("	}")
	gFile.Println("")
}

func generateFutureTableRowStruct(gFile *GeneratedFile, table *scheme.Table) {
	gFile.Println("type Future" + table.Name + "TableRow struct {")
	gFile.Println("	Dec lib.Decoder")
	gFile.Println("	Future fdb.FutureByteSlice")
	gFile.Println("}")
	gFile.Println("")
}

func generateMethods(gFile *GeneratedFile, table *scheme.Table, models []*scheme.Model) {
	generateTableMethods(gFile, table)
	generateTablePKMethods(gFile, table, models)
	generateFutureTableRowMethods(gFile, table)
}

func generateTableMethods(gFile *GeneratedFile, table *scheme.Table) {
	gFile.Println("func (table *" + table.Name + "Table) Get(tr fdb.ReadTransaction, pk *" + table.Name + "TablePK) (*Future" + table.Name + "TableRow, error) {")
	gFile.Println("	key, err := pk.Pack()")
	gFile.Println("	if err != nil {")
	gFile.Println("		return nil, err")
	gFile.Println("	}")
	gFile.Println("")
	gFile.Println("	future := tr.Get(table.Sub.Sub(key...))")
	gFile.Println("	return &Future" + table.Name + "TableRow{Future: future, Dec: table.Dec}, nil")
	gFile.Println("}")
	gFile.Println("")

	gFile.Println("func (table *" + table.Name + "Table) MustGet(tr fdb.ReadTransaction, pk *" + table.Name + "TablePK) *Future" + table.Name + "TableRow {")
	gFile.Println("	future, err := table.Get(tr, pk)")
	gFile.Println("	if err != nil {")
	gFile.Println("		panic(err)")
	gFile.Println("	}")
	gFile.Println("	return future")
	gFile.Println("}")
	gFile.Println("")

	gFile.Println("func (table *" + table.Name + "Table) Insert(tr fdb.Transaction, model *" + table.Name + "TableRow) error {")
	gFile.Println("	pk := &" + table.Name + "TablePK{")
	for _, column := range table.PK {
		gFile.Println("		" + strings.Join(strings.Split(column, "."), "") + ": model." + column + ",")
	}
	gFile.Println("	}")
	gFile.Println("	key, err := pk.Pack()")
	gFile.Println("	if err != nil {")
	gFile.Println("		return err")
	gFile.Println("	}")
	gFile.Println("")
	gFile.Println("	value, err := table.Enc.Encode(model)")
	gFile.Println("	if err != nil {")
	gFile.Println("		return err")
	gFile.Println("	}")
	gFile.Println("	tr.Set(table.Sub.Sub(key...), value)")
	gFile.Println("	return nil")
	gFile.Println("}")
	gFile.Println("")

	gFile.Println("func (table *" + table.Name + "Table) MustInsert(tr fdb.Transaction, model *" + table.Name + "TableRow) {")
	gFile.Println("	err := table.Insert(tr, model)")
	gFile.Println("	if err != nil {")
	gFile.Println("		panic(err)")
	gFile.Println("	}")
	gFile.Println("}")
	gFile.Println("")

	gFile.Println("func (table *" + table.Name + "Table) Delete(tr fdb.Transaction, pk *" + table.Name + "TablePK) error {")
	gFile.Println("	key, err := pk.Pack()")
	gFile.Println("	if err != nil {")
	gFile.Println("		return err")
	gFile.Println("	}")
	gFile.Println("	tr.Clear(table.Sub.Sub(key...))")
	gFile.Println("	return nil")
	gFile.Println("}")
	gFile.Println("")

	gFile.Println("func (table *" + table.Name + "Table) MustDelete(tr fdb.Transaction, pk *" + table.Name + "TablePK) {")
	gFile.Println("	err := table.Delete(tr, pk)")
	gFile.Println("	if err != nil {")
	gFile.Println("		panic(err)")
	gFile.Println("	}")
	gFile.Println("}")
	gFile.Println("")
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
			gFile.Println("	pk" + name + "Bytes := []byte(pk." + name + ")")
		case "uint64":
			gFile.Println("	" + name + "Buf := new(bytes.Buffer)")
			gFile.Println("	err = binary.Write(" + name + "Buf, binary.BigEndian, pk." + name + ")")
			gFile.Println("	if err != nil {")
			gFile.Println("		return nil, err")
			gFile.Println("	}")
			gFile.Println("	pk" + name + "Bytes := " + name + "Buf.Bytes()")
		case "int64":
			gFile.Println("	" + name + "Buf := new(bytes.Buffer)")
			gFile.Println("	err = binary.Write(" + name + "Buf, binary.BigEndian, " + name + ")")
			gFile.Println("	if err != nil {")
			gFile.Println("		return nil, err")
			gFile.Println("	}")
			gFile.Println("	pk" + name + "Bytes := " + name + "Buf.Bytes()")
		case "float32":
			gFile.Println("	" + name + "Buf := new(bytes.Buffer)")
			gFile.Println("	err = binary.Write(" + name + "Buf, binary.BigEndian, " + name + ")")
			gFile.Println("	if err != nil {")
			gFile.Println("		return nil, err")
			gFile.Println("	}")
			gFile.Println("	pk" + name + "Bytes := " + name + "Buf.Bytes()")
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
	gFile.Println("")
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
	gFile.Println("")

	funcGetString :=
		`func (future *Future%[1]sTableRow) Get() (*%[1]sTableRow, error) {
		value, err := future.Future.Get()
		if err != nil {
			return nil, err
		}
		return future.new%[1]sTableRow(value)
	}`
	gFile.Println(fmt.Sprintf(funcGetString, table.Name))
	gFile.Println("")

	funcMustGetString :=
		`func (future *Future%[1]sTableRow) MustGet() *%[1]sTableRow {
		value, err := future.Get()
		if err != nil {
			panic(err)
		}
		return value
	}`
	gFile.Println(fmt.Sprintf(funcMustGetString, table.Name))
	gFile.Println("")
}
