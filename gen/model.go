package gen

import (
	"fmt"
	"github.com/nostressdev/fdb/orm/scheme"
)

func GenerateModel(gFile *GeneratedFile, model *scheme.Model) {
	gFile.Println(fmt.Sprintf("type %s struct {", model.Name))
	for _, field := range model.Fields {
		gFile.Println(fmt.Sprintf("	%s %s", field.Name, field.Type))
	}
	gFile.Println("}")
	gFile.Println("")
}
