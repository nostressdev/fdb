package lib

type Encoder interface {
	Encode(interface{}) ([]byte, error)
}

type Decoder interface {
	Decode([]byte) (interface{}, error)
}
