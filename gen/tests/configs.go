package tests

import "github.com/nostressdev/fdb/orm/scheme"

var Config = &scheme.GeneratorConfig{
	FilesPath:   "generated/",
	PackageName: "gen_fdb",
	Models: []*scheme.Model{{
		Name: "User",
		Fields: []*scheme.Field{
			{
				Name: "ID",
				Type: "string",
			},
			{
				Name: "Age",
				Type: "bool",
			},
		},
	}},
	Tables: []*scheme.Table{
		{
			Name:    "Users",
			Columns: []*scheme.Column{{Name: "Man", Type: "User"}, {Name: "Ts", Type: "uint64"}},
			PK:      []string{"Ts", "Man.ID"},
		},
		{
			Name:    "AgeSort",
			Columns: []*scheme.Column{{Name: "Man", Type: "User"}},
			PK:      []string{"Man.Age", "Man.ID"},
		},
	},
}
