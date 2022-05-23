package gen

import (
	"fmt"
	"github.com/nostressdev/fdb/orm/scheme"
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
)

func GenerateTable(gFile *protogen.GeneratedFile, table *scheme.Table, models []*scheme.Model) {
	generateStructs(gFile, table, models)
	generateMethods(gFile, table, models)
}

func generateStructs(gFile *protogen.GeneratedFile, table *scheme.Table, models []*scheme.Model) {
	generateTableStruct(gFile, table)
	generateTablePKStruct(gFile, table, models)
	generateTableRowStruct(gFile, table)
	generateFutureTableRowStruct(gFile, table)
}

func generateTableStruct(gFile *protogen.GeneratedFile, table *scheme.Table) {
	gFile.P("type " + table.Name + "Table struct {")
	gFile.P("	Enc lib.Encoder")
	gFile.P("	Dec lib.Decoder")
	gFile.P("	Sub subspace.Subspace")
	gFile.P("}")
	gFile.P()
	gFile.P("func New" + table.Name + "Table(opts ...lib.TableOptions) (*" + table.Name + "Table, error) {")
	gFile.P("	table := &" + table.Name + "Table{}")
	gFile.P("	for _, opt := range opts {")
	gFile.P("		if opt.Enc != nil {")
	gFile.P("			table.Enc = opt.Enc")
	gFile.P("		}")
	gFile.P("		if opt.Dec != nil {")
	gFile.P("			table.Dec = opt.Dec")
	gFile.P("		}")
	gFile.P("		if opt.Sub != nil {")
	gFile.P("			table.Sub = opt.Sub")
	gFile.P("		}")
	gFile.P("	}")
	gFile.P("	if table.Enc == nil {")
	gFile.P("		return nil, fmt.Errorf(\"encoder is nil\")")
	gFile.P("	}")
	gFile.P("	if table.Dec == nil {")
	gFile.P("		return nil, fmt.Errorf(\"decoder is nil\")")
	gFile.P("	}")
	gFile.P("	if table.Sub == nil {")
	gFile.P("		return nil, fmt.Errorf(\"subspace is nil\")")
	gFile.P("	}")
	gFile.P("	return table, nil")
	gFile.P("}")
	gFile.P()
}

func generateTablePKStruct(gFile *protogen.GeneratedFile, table *scheme.Table, models []*scheme.Model) {
	gFile.P("type " + table.Name + "TablePK struct {")
	for _, column := range table.PK {
		tp, err := getType(table, models, column)
		if err != nil {
			panic(err)
		}
		gFile.P("		" + strings.Join(strings.Split(column, "."), "") + " " + tp)
	}
	gFile.P("	}")
	gFile.P()
}

func generateTableRowStruct(gFile *protogen.GeneratedFile, table *scheme.Table) {
	gFile.P("type " + table.Name + "TableRow struct {")
	for _, column := range table.Columns {
		gFile.P("		" + column.Name + " " + column.Type)
	}
	gFile.P("	}")
	gFile.P()
}

func generateFutureTableRowStruct(gFile *protogen.GeneratedFile, table *scheme.Table) {
	gFile.P("type Future" + table.Name + "TableRow struct {")
	gFile.P("	Dec lib.Decoder")
	gFile.P("	Future fdb.FutureByteSlice")
	gFile.P("}")
	gFile.P()
}

func generateMethods(gFile *protogen.GeneratedFile, table *scheme.Table, models []*scheme.Model) {
	generateTableMethods(gFile, table)
	generateTablePKMethods(gFile, table, models)
	generateFutureTableRowMethods(gFile, table)
}

func generateTableMethods(gFile *protogen.GeneratedFile, table *scheme.Table) {
	gFile.P("func (table *" + table.Name + "Table) Get(tr fdb.ReadTransaction, pk *" + table.Name + "TablePK) (*Future" + table.Name + "TableRow, error) {")
	gFile.P("	key, err := pk.Pack()")
	gFile.P("	if err != nil {")
	gFile.P("		return nil, err")
	gFile.P("	}")
	gFile.P()
	gFile.P("	fmt.Println(table.Sub.Sub(key...).Bytes())")
	gFile.P("	future := tr.Get(table.Sub.Sub(key...))")
	gFile.P("	return &Future" + table.Name + "TableRow{Future: future, Dec: table.Dec}, nil")
	gFile.P("}")
	gFile.P()

	gFile.P("func (table *" + table.Name + "Table) MustGet(tr fdb.ReadTransaction, pk *" + table.Name + "TablePK) *Future" + table.Name + "TableRow {")
	gFile.P("	future, err := table.Get(tr, pk)")
	gFile.P("	if err != nil {")
	gFile.P("		panic(err)")
	gFile.P("	}")
	gFile.P("	return future")
	gFile.P("}")
	gFile.P()

	gFile.P("func (table *" + table.Name + "Table) Insert(tr fdb.Transaction, model *" + table.Name + "TableRow) error {")
	gFile.P("	pk := &" + table.Name + "TablePK{")
	for _, column := range table.PK {
		gFile.P("		" + strings.Join(strings.Split(column, "."), "") + ": model." + column + ",")
	}
	gFile.P("	}")
	gFile.P("	key, err := pk.Pack()")
	gFile.P("	if err != nil {")
	gFile.P("		return err")
	gFile.P("	}")
	gFile.P()
	gFile.P("	value, err := table.Enc.Encode(model)")
	gFile.P("	if err != nil {")
	gFile.P("		return err")
	gFile.P("	}")
	gFile.P("	fmt.Println(string(value))")
	gFile.P("	fmt.Println(table.Sub.Sub(key...).Bytes())")
	gFile.P("	tr.Set(table.Sub.Sub(key...), value)")
	gFile.P("	return nil")
	gFile.P("}")
	gFile.P()

	gFile.P("func (table *" + table.Name + "Table) MustInsert(tr fdb.Transaction, model *" + table.Name + "TableRow) {")
	gFile.P("	err := table.Insert(tr, model)")
	gFile.P("	if err != nil {")
	gFile.P("		panic(err)")
	gFile.P("	}")
	gFile.P("}")
	gFile.P()

	gFile.P("func (table *" + table.Name + "Table) Delete(tr fdb.Transaction, pk *" + table.Name + "TablePK) error {")
	gFile.P("	key, err := pk.Pack()")
	gFile.P("	if err != nil {")
	gFile.P("		return err")
	gFile.P("	}")
	gFile.P("	tr.Clear(table.Sub.Sub(key...))")
	gFile.P("	return nil")
	gFile.P("}")
	gFile.P()

	gFile.P("func (table *" + table.Name + "Table) MustDelete(tr fdb.Transaction, pk *" + table.Name + "TablePK) {")
	gFile.P("	err := table.Delete(tr, pk)")
	gFile.P("	if err != nil {")
	gFile.P("		panic(err)")
	gFile.P("	}")
	gFile.P("}")
	gFile.P()
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

func generateTablePKMethods(gFile *protogen.GeneratedFile, table *scheme.Table, models []*scheme.Model) {
	gFile.P("func (pk *" + table.Name + "TablePK) Pack() ([]tuple.TupleElement, error) {")
	gFile.P("	var err error")
	elements := make([]string, 0, len(table.PK))
	for _, column := range table.PK {
		name := strings.Join(strings.Split(column, "."), "")
		tp, err := getType(table, models, column)
		if err != nil {
			panic(err)
		}
		switch tp {
		case "string":
			gFile.P("	pk" + name + "Bytes := []byte(pk." + name + ")")
		case "uint64":
			gFile.P("	" + name + "Buf := new(bytes.Buffer)")
			gFile.P("	err = binary.Write(" + name + "Buf, binary.BigEndian, pk." + name + ")")
			gFile.P("	if err != nil {")
			gFile.P("		return nil, err")
			gFile.P("	}")
			gFile.P("	pk" + name + "Bytes := " + name + "Buf.Bytes()")
		case "int64":
			gFile.P("	" + name + "Buf := new(bytes.Buffer)")
			gFile.P("	err = binary.Write(" + name + "Buf, binary.BigEndian, " + name + ")")
			gFile.P("	if err != nil {")
			gFile.P("		return nil, err")
			gFile.P("	}")
			gFile.P("	pk" + name + "Bytes := " + name + "Buf.Bytes()")
		case "float32":
			gFile.P("	" + name + "Buf := new(bytes.Buffer)")
			gFile.P("	err = binary.Write(" + name + "Buf, binary.BigEndian, " + name + ")")
			gFile.P("	if err != nil {")
			gFile.P("		return nil, err")
			gFile.P("	}")
			gFile.P("	pk" + name + "Bytes := " + name + "Buf.Bytes()")
		default:
			panic("unknown type " + tp)
		}
		elements = append(elements, "pk"+name+"Bytes")
	}
	gFile.P("	if err != nil {")
	gFile.P("		return nil, err")
	gFile.P("	}")
	gFile.P("	return []tuple.TupleElement{" + strings.Join(elements, ", ") + "}, nil")
	gFile.P("}")
	gFile.P()
}

func generateFutureTableRowMethods(gFile *protogen.GeneratedFile, table *scheme.Table) {
	funcNewTableRowString :=
		`func (future *Future%[1]sTableRow) new%[1]sTableRow(value []byte) (*%[1]sTableRow, error) {
			if row, err := future.Dec.Decode(value); err != nil {
				return nil, err
			} else {
				return row.(*%[1]sTableRow), nil
			}
		}`
	gFile.P(fmt.Sprintf(funcNewTableRowString, table.Name))
	gFile.P()

	funcGetString :=
		`func (future *Future%[1]sTableRow) Get() (*%[1]sTableRow, error) {
		value, err := future.Future.Get()
		if err != nil {
			return nil, err
		}
		fmt.Println(string(value))
		return future.new%[1]sTableRow(value)
	}`
	gFile.P(fmt.Sprintf(funcGetString, table.Name))
	gFile.P()

	funcMustGetString :=
		`func (future *Future%[1]sTableRow) MustGet() *%[1]sTableRow {
		value, err := future.Get()
		if err != nil {
			panic(err)
		}
		return value
	}`
	gFile.P(fmt.Sprintf(funcMustGetString, table.Name))
	gFile.P()
}
