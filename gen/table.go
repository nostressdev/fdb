package gen

import (
	"fmt"
	"github.com/nostressdev/fdb/orm/scheme"
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
)

func GenerateTable(gFile *protogen.GeneratedFile, table *scheme.Table, models []*scheme.Model) {
	generateStructs(gFile, table)
	generateMethods(gFile, table, models)
}

func generateStructs(gFile *protogen.GeneratedFile, table *scheme.Table) {
	generateTableStruct(gFile, table)
	generateTablePKStruct(gFile, table)
	generateTableRowStruct(gFile, table)
	generateFutureTableRowStruct(gFile, table)
}

func generateTableStruct(gFile *protogen.GeneratedFile, table *scheme.Table) {
	gFile.P("type " + table.Name + "Table struct {")
	gFile.P("	Enc lib.Encoder")
	gFile.P("	Sub subspace.Subspace")
	gFile.P("}")
	gFile.P()
	gFile.P("func New" + table.Name + "Table(opts ...lib.TableOptions) (*" + table.Name + "Table, error) {")
	gFile.P("	table := &" + table.Name + "Table{}")
	gFile.P("	for _, opt := range opts {")
	gFile.P("		if opt.Enc != nil {")
	gFile.P("			table.Enc = opt.Enc")
	gFile.P("		}")
	gFile.P("		if opt.Sub != nil {")
	gFile.P("			table.Sub = opt.Sub")
	gFile.P("		}")
	gFile.P("	}")
	gFile.P("	if table.Enc == nil {")
	gFile.P("		return nil, fmt.Errorf(\"encoder is nil\")")
	gFile.P("	}")
	gFile.P("	if table.Sub == nil {")
	gFile.P("		return nil, fmt.Errorf(\"subspace is nil\")")
	gFile.P("	}")
	gFile.P("	return table, nil")
	gFile.P("}")
	gFile.P()
}

func generateTablePKStruct(gFile *protogen.GeneratedFile, table *scheme.Table) {
	gFile.P("type " + table.Name + "TablePK struct {")
	for _, column := range table.PK {
		gFile.P("		" + strings.Join(strings.Split(column, "."), "") + " lib.Key")
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
	gFile.P("	Enc lib.Encoder")
	gFile.P("	Future fdb.FutureByteSlice")
	gFile.P("}")
	gFile.P()
}

func generateMethods(gFile *protogen.GeneratedFile, table *scheme.Table, models []*scheme.Model) {
	generateTableMethods(gFile, table, models)
	generateTablePKMethods(gFile, table)
	generateFutureTableRowMethods(gFile, table)
}

func generateTableMethods(gFile *protogen.GeneratedFile, table *scheme.Table, models []*scheme.Model) {
	gFile.P("func (table *" + table.Name + "Table) Get(tr fdb.ReadTransaction, pk *" + table.Name + "TablePK) (*Future" + table.Name + "TableRow, error) {")
	gFile.P("	key, err := pk.Pack()")
	gFile.P("	if err != nil {")
	gFile.P("		return nil, err")
	gFile.P("	}")
	gFile.P()
	gFile.P("	future := tr.Get(table.Sub.Sub(key))")
	gFile.P("	return &Future" + table.Name + "TableRow{Future: future, Enc: table.Enc}, nil")
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
		tp, err := getType(table, models, column)
		if err != nil {
			panic(err)
		}
		gFile.P("		" + strings.Join(strings.Split(column, "."), "") + ": &lib.Key" + strings.Title(tp) + "{Value: model." + column + "},")
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
	gFile.P("	tr.Set(table.Sub.Sub(key), value)")
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
	gFile.P("	tr.Clear(table.Sub.Sub(key))")
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
	qSlice := strings.Split(q, ".")
	for _, column := range table.Columns {
		if column.Name == qSlice[0] {
			for _, model := range models {
				if column.Type == model.Name {
					for _, field := range model.Fields {
						if field.Name == qSlice[1] {
							return field.Type, nil
						}
					}
					return "", fmt.Errorf("don't find \"%s\" from \"%s\" in model \"%s\"", qSlice[1], q, model.Name)
				}
			}
			return "", fmt.Errorf("don't find type \"%s\" of column \"%s\" from \"%s\"", column.Type, column.Name, q)
		}
	}
	return "", fmt.Errorf("don't find \"%s\" from \"%s\" in table columns", qSlice[0], q)
}

func generateTablePKMethods(gFile *protogen.GeneratedFile, table *scheme.Table) {
	gFile.P("func (pk *" + table.Name + "TablePK) Pack() ([]tuple.TupleElement, error) {")
	elements := make([]string, 0, len(table.PK))
	for _, column := range table.PK {
		gFile.P("	pk" + strings.Join(strings.Split(column, "."), "") + "Bytes, err := pk." + strings.Join(strings.Split(column, "."), "") + ".Key()")
		gFile.P("	if err != nil {")
		gFile.P("		return nil, err")
		gFile.P("	}")
		elements = append(elements, "pk"+strings.Join(strings.Split(column, "."), "")+"Bytes")
	}
	gFile.P("	return []tuple.TupleElement{" + strings.Join(elements, ", ") + "}, nil")
	gFile.P("}")
	gFile.P()
}

func generateFutureTableRowMethods(gFile *protogen.GeneratedFile, table *scheme.Table) {
	gFile.P("func (future *Future" + table.Name + "TableRow) new" + table.Name + "TableRow(value []byte) (*" + table.Name + "TableRow, error) {")
	gFile.P("	if row, err := future.Enc.Decode(value); err != nil {")
	gFile.P("		return nil, err")
	gFile.P("	} else {")
	gFile.P("		return row.(*" + table.Name + "TableRow), nil")
	gFile.P("	}")
	gFile.P("}")
	gFile.P()

	gFile.P("func (future *Future" + table.Name + "TableRow) Get() (*" + table.Name + "TableRow, error) {")
	gFile.P("	value, err := future.Future.Get()")
	gFile.P("	if err != nil {")
	gFile.P("		return nil, err")
	gFile.P("	}")
	gFile.P("	return future.new" + table.Name + "TableRow(value)")
	gFile.P("}")
	gFile.P()

	gFile.P("func (future *Future" + table.Name + "TableRow) MustGet() *" + table.Name + "TableRow {")
	gFile.P("	value, err := future.Get()")
	gFile.P("	if err != nil {")
	gFile.P("		panic(err)")
	gFile.P("	}")
	gFile.P("	return value")
	gFile.P("}")
	gFile.P()
}
