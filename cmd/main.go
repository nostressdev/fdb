package main

import (
	"github.com/nostressdev/fdb/gen"
	"github.com/nostressdev/fdb/orm/scheme"
)

func main() {
	config := &scheme.GeneratorConfig{}
	config.FilesPath = "generated/"
	config.PackageName = "gen"
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
	config.Tables = []*scheme.Table{
		{
			Name:    "Users",
			Columns: []scheme.Column{{Name: "Man", Type: "User"}, {Name: "Ts", Type: "uint64"}},
			PK:      []string{"Ts", "Man.ID"},
		},
		{
			Name:    "AgeSort",
			Columns: []scheme.Column{{Name: "Man", Type: "User"}},
			PK:      []string{"Man.Age", "Man.ID"},
		},
	}

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

	gen.GenFiles(config)
}
