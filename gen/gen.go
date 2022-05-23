package gen

import (
	"github.com/nostressdev/fdb/orm/scheme"
	"go/format"
	"google.golang.org/protobuf/compiler/protogen"
	"os"
)

func GenFiles(config *scheme.GeneratorConfig) {
	err := os.RemoveAll(config.FilesPath)
	if err != nil {
		panic(err)
	}

	err = os.Mkdir(config.FilesPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	gFile := &protogen.GeneratedFile{}
	GenModels(gFile, config)
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
		GenTable(gFile, config, i)
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

		gFile = &protogen.GeneratedFile{}
		GenEncoder(gFile, config, i)
		c, err = gFile.Content()
		if err != nil {
			panic(err)
		}
		c, err = format.Source(c)
		if err != nil {
			panic(err)
		}
		f, err = os.Create(config.FilesPath + table.Name + "TableEncoder.g.go")
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

func GenModels(gFile *protogen.GeneratedFile, config *scheme.GeneratorConfig) {
	gFile.P("package " + config.PackageName)
	gFile.P()
	for _, model := range config.Models {
		GenerateModel(gFile, model)
	}
}

func GenTable(gFile *protogen.GeneratedFile, config *scheme.GeneratorConfig, index int) {
	gFile.P("package " + config.PackageName)
	gFile.P()
	gFile.P("import (")
	gFile.P("	\"bytes\"")
	gFile.P("	\"encoding/binary\"")
	gFile.P("	\"fmt\"")
	gFile.P("	\"github.com/apple/foundationdb/bindings/go/src/fdb\"")
	gFile.P("	\"github.com/apple/foundationdb/bindings/go/src/fdb/subspace\"")
	gFile.P("	\"github.com/apple/foundationdb/bindings/go/src/fdb/tuple\"")
	gFile.P("	\"github.com/nostressdev/fdb/lib\"")
	gFile.P(")")
	gFile.P()
	GenerateTable(gFile, config.Tables[index], config.Models)
}

func GenEncoder(gFile *protogen.GeneratedFile, config *scheme.GeneratorConfig, index int) {
	gFile.P("package " + config.PackageName)
	gFile.P()
	gFile.P("import (")
	gFile.P("	\"encoding/json\"")
	gFile.P(")")
	gFile.P()
	GenerateEncoder(gFile, config.Tables[index])
	GenerateDecoder(gFile, config.Tables[index])
}
