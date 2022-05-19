package main

import (
	"github.com/nostressdev/fdb/gen"
	"github.com/nostressdev/fdb/orm/scheme"
	"go/format"
	"google.golang.org/protobuf/compiler/protogen"
	"os"
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

	err := os.RemoveAll(config.FilesPath)
	if err != nil {
		panic(err)
	}

	err = os.Mkdir(config.FilesPath, 7770)
	if err != nil {
		panic(err)
	}

	gFile := &protogen.GeneratedFile{}
	gen.GenModels(gFile, config)
	c, err := gFile.Content()
	if err != nil {
		panic(err)
	}
	c, err = format.Source(c)
	if err != nil {
		panic(err)
	}
	f, err := os.Create(config.FilesPath + "Models.g.go")
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
	for i, table := range config.Tables {
		gFile = &protogen.GeneratedFile{}
		gen.GenTable(gFile, config, i)
		c, err = gFile.Content()
		if err != nil {
			panic(err)
		}
		c, err = format.Source(c)
		if err != nil {
			panic(err)
		}
		f, err = os.Create(config.FilesPath + table.Name + "Table.g.go")
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
}
