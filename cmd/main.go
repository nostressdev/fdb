package main

import (
	"github.com/nostressdev/fdb/gen"
	"github.com/nostressdev/fdb/orm/scheme"
	"google.golang.org/protobuf/compiler/protogen"
	"os"
)

func main() {
	config := &scheme.GeneratorConfig{}
	config.Models = []*scheme.Model{{
		Name: "User",
		Fields: []scheme.Field{
			{
				Name: "ID",
				Type: "string",
			},
			{
				Name: "Age",
				Type: "uint64",
			},
		},
	}}
	config.Tables = []*scheme.Table{{
		Name:    "Users",
		Columns: []scheme.Column{{Name: "Man", Type: "User"}, {Name: "Ts", Type: "uint64"}},
		PK:      []string{"Ts", "Man.ID"},
	}}

	//protogen.Options{}.Run(func(gen *protogen.Plugin) error {
	//	fmt.Println(0)
	//	for _, f := range gen.Files {
	//		if !f.Generate {
	//			continue
	//		}
	//		generateFile(gen, f)
	//	}
	//	return nil
	//})

	gFile := &protogen.GeneratedFile{}
	gen.GenerateCode(gFile, config)
	c, err := gFile.Content()
	if err != nil {
		panic(err)
	}
	f, err := os.Create("genFile.go")
	if err != nil {
		panic(err)
	}
	_, err = f.Write(c)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}
}
