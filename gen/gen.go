package gen

import (
	"github.com/nostressdev/fdb/orm/scheme"
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

	gFile := &GeneratedFile{}
	GenModels(gFile, config)
	f, err := os.Create(config.FilesPath + "Models.g.go")
	if err != nil {
		panic(err)
	}
	gFile.Write(f)
	err = f.Close()
	if err != nil {
		panic(err)
	}
	for i, table := range config.Tables {
		gFile = &GeneratedFile{}
		GenTable(gFile, config, i)
		f, err = os.Create(config.FilesPath + table.Name + "Table.g.go")
		if err != nil {
			panic(err)
		}
		gFile.Write(f)
		err = f.Close()
		if err != nil {
			panic(err)
		}

		gFile = &GeneratedFile{}
		GenEncoder(gFile, config, i)
		f, err = os.Create(config.FilesPath + table.Name + "TableEncoder.g.go")
		if err != nil {
			panic(err)
		}
		gFile.Write(f)
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}
}

func GenModels(gFile *GeneratedFile, config *scheme.GeneratorConfig) {
	gFile.Println("package " + config.PackageName)
	gFile.Println("")
	for _, model := range config.Models {
		GenerateModel(gFile, model)
	}
}

func GenTable(gFile *GeneratedFile, config *scheme.GeneratorConfig, index int) {
	gFile.Println("package " + config.PackageName)
	gFile.Println("")
	gFile.Println("import (")
	gFile.Println("	\"bytes\"")
	gFile.Println("	\"encoding/binary\"")
	gFile.Println("	\"fmt\"")
	gFile.Println("	\"github.com/apple/foundationdb/bindings/go/src/fdb\"")
	gFile.Println("	\"github.com/apple/foundationdb/bindings/go/src/fdb/subspace\"")
	gFile.Println("	\"github.com/apple/foundationdb/bindings/go/src/fdb/tuple\"")
	gFile.Println("	\"github.com/nostressdev/fdb/lib\"")
	gFile.Println(")")
	gFile.Println("")
	GenerateTable(gFile, config.Tables[index], config.Models)
}

func GenEncoder(gFile *GeneratedFile, config *scheme.GeneratorConfig, index int) {
	gFile.Println("package " + config.PackageName)
	gFile.Println("")
	gFile.Println("import (")
	gFile.Println("	\"encoding/json\"")
	gFile.Println(")")
	gFile.Println("")
	GenerateEncoder(gFile, config.Tables[index])
	GenerateDecoder(gFile, config.Tables[index])
}
