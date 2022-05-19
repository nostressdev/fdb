package lib

import "encoding/json"

type Encoder interface {
	Encode(interface{}) ([]byte, error)
	Decode([]byte) (interface{}, error)
}

type JsonEncoder struct{}

func (enc *JsonEncoder) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (enc *JsonEncoder) Decode(value []byte) (interface{}, error) {
	var res interface{}
	err := json.Unmarshal(value, res)
	return res, err
}
