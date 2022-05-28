package gen

import (
	"github.com/nostressdev/fdb/orm/scheme"
)

func GenerateModel(gFile *GeneratedFile, model *scheme.Model) {
	gFile.Printf("type %s struct {", model.Name)
	for _, field := range model.Fields {
		gFile.Printf("	%s %s", field.Name, field.Type)
	}
	gFile.Println("}")
	gFile.Println("")
}
