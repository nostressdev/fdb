package gen

import (
	"github.com/nostressdev/fdb/orm/scheme"
)

func GenerateModel(gFile *GeneratedFile, model *scheme.Model) {
	gFile.Println("type " + model.Name + " struct {")
	for _, field := range model.Fields {
		gFile.Println(" 	" + field.Name + " " + field.Type)
	}
	gFile.Println("}")
	gFile.Println("")
}
