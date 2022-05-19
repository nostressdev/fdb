package gen

import (
	"github.com/nostressdev/fdb/orm/scheme"
	"google.golang.org/protobuf/compiler/protogen"
)

func GenerateModel(gFile *protogen.GeneratedFile, model *scheme.Model) {
	gFile.P("type " + model.Name + " struct {")
	for _, field := range model.Fields {
		gFile.P(" 	" + field.Name + " " + field.Type)
	}
	gFile.P("}")
	gFile.P()
}
