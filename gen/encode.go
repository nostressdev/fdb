package gen

import (
	"github.com/nostressdev/fdb/orm/scheme"
)

func GenerateEncoder(gFile *GeneratedFile, table *scheme.Table) {
	gFile.Println("type " + table.Name + "TableRowJsonEncoder struct {")
	gFile.Println("}")
	gFile.Println("")
	gFile.Println("func (enc *" + table.Name + "TableRowJsonEncoder) Encode(value interface{}) ([]byte, error){")
	gFile.Println("	res, err := json.Marshal(value)")
	gFile.Println("	if err != nil {")
	gFile.Println("		panic(err)")
	gFile.Println("	}")
	gFile.Println("	return res, nil")
	gFile.Println("}")
}

func GenerateDecoder(gFile *GeneratedFile, table *scheme.Table) {
	gFile.Println("type " + table.Name + "TableRowJsonDecoder struct {")
	gFile.Println("}")
	gFile.Println("")
	gFile.Println("func (enc *" + table.Name + "TableRowJsonDecoder) Decode(value []byte) (interface{}, error){")
	gFile.Println("	res := &" + table.Name + "TableRow{}")
	gFile.Println("	err := json.Unmarshal(value, res)")
	gFile.Println("	if err != nil {")
	gFile.Println("		panic(err)")
	gFile.Println("	}")
	gFile.Println("	return res, nil")
	gFile.Println("}")
}
