package gen

import (
	"github.com/nostressdev/fdb/orm/scheme"
)

func GenerateEncoder(gFile *GeneratedFile, table *scheme.Table) {
	tableRowJsonEncoderString := `type %sTableRowJsonEncoder struct {}`
	encodeString :=
		`func (enc *%sTableRowJsonEncoder) Encode(value interface{}) ([]byte, error){
			res, err := json.Marshal(value)
			if err != nil {
				panic(err)
			}
			return res, nil
		}`

	gFile.Printf(tableRowJsonEncoderString, table.Name)
	gFile.Printf(encodeString, table.Name)
}

func GenerateDecoder(gFile *GeneratedFile, table *scheme.Table) {
	tableRowJsonDecoderString := `type %sTableRowJsonDecoder struct {}`
	decodeString :=
		`func (dec *%[1]sTableRowJsonDecoder) Decode(value []byte) (interface{}, error){
			res := &%[1]sTableRow{}
			err := json.Unmarshal(value, res)
			if err != nil {
				panic(err)
			}
			return res, nil
		}`

	gFile.Printf(tableRowJsonDecoderString, table.Name)
	gFile.Printf(decodeString, table.Name)
}
