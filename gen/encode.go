package gen

import (
	"github.com/nostressdev/fdb/orm/scheme"
	"google.golang.org/protobuf/compiler/protogen"
)

func GenerateEncoder(gFile *protogen.GeneratedFile, table *scheme.Table) {
	gFile.P("type " + table.Name + "TableRowJsonEncoder struct {")
	gFile.P("}")
	gFile.P()
	gFile.P("func (enc *" + table.Name + "TableRowJsonEncoder) Encode(value interface{}) ([]byte, error){")
	gFile.P("	res, err := json.Marshal(value)")
	gFile.P("	if err != nil {")
	gFile.P("		panic(err)")
	gFile.P("	}")
	gFile.P("	return res, nil")
	gFile.P("}")
}

func GenerateDecoder(gFile *protogen.GeneratedFile, table *scheme.Table) {
	gFile.P("type " + table.Name + "TableRowJsonDecoder struct {")
	gFile.P("}")
	gFile.P()
	gFile.P("func (enc *" + table.Name + "TableRowJsonDecoder) Decode(value []byte) (interface{}, error){")
	gFile.P("	res := &" + table.Name + "TableRow{}")
	gFile.P("	err := json.Unmarshal(value, res)")
	gFile.P("	if err != nil {")
	gFile.P("		panic(err)")
	gFile.P("	}")
	gFile.P("	return res, nil")
	gFile.P("}")
}
