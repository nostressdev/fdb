package gen

import (
	"github.com/nostressdev/fdb/orm/scheme"
	"google.golang.org/protobuf/compiler/protogen"
)

func GenerateCode(gFile *protogen.GeneratedFile, config *scheme.GeneratorConfig) {
	gFile.P("package genFdb")
	gFile.P()
	gFile.P("import (")
	gFile.P("	\"fmt\"")
	gFile.P("	\"github.com/apple/foundationdb/bindings/go/src/fdb\"")
	gFile.P("	\"github.com/apple/foundationdb/bindings/go/src/fdb/subspace\"")
	gFile.P("	\"github.com/apple/foundationdb/bindings/go/src/fdb/tuple\"")
	gFile.P("	\"github.com/nostressdev/fdb/lib\"")
	gFile.P(")")
	gFile.P()
	for _, model := range config.Models {
		GenerateModel(gFile, model)
	}
	for _, table := range config.Tables {
		GenerateTable(gFile, table, config.Models)
	}
}

func GenModels(gFile *protogen.GeneratedFile, config *scheme.GeneratorConfig) {
	gFile.P("package genFdb")
	gFile.P()
	for _, model := range config.Models {
		GenerateModel(gFile, model)
	}
}

func GenTable(gFile *protogen.GeneratedFile, config *scheme.GeneratorConfig, index int) {
	gFile.P("package genFdb")
	gFile.P()
	gFile.P("import (")
	gFile.P("	\"fmt\"")
	gFile.P("	\"github.com/apple/foundationdb/bindings/go/src/fdb\"")
	gFile.P("	\"github.com/apple/foundationdb/bindings/go/src/fdb/subspace\"")
	gFile.P("	\"github.com/apple/foundationdb/bindings/go/src/fdb/tuple\"")
	gFile.P("	\"github.com/nostressdev/fdb/lib\"")
	gFile.P(")")
	gFile.P()
	GenerateTable(gFile, config.Tables[index], config.Models)
}
