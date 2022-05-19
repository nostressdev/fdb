package lib

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"sync"
)

var bytesBufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func AcquireBytesBuffer() *bytes.Buffer {
	buf := bytesBufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

func ReleaseBytesBuffer(bytes *bytes.Buffer) {
	bytes.Reset()
	bytesBufferPool.Put(bytes)
}

func GetBigEndianBytesUint64(n uint64) (*bytes.Buffer, error) {
	buf := AcquireBytesBuffer()
	err := binary.Write(buf, binary.BigEndian, n)
	return buf, err
}

func MustGetBigEndianBytesUint64(n uint64) *bytes.Buffer {
	buf := AcquireBytesBuffer()
	err := binary.Write(buf, binary.BigEndian, n)
	if err != nil {
		panic(err)
	}
	return buf
}

type Key interface {
	Key() ([]byte, error)
	MustKey() []byte
}

type KeyString struct {
	Value string
}

func (key *KeyString) Key() ([]byte, error) {
	return []byte(key.Value), nil
}

func (key *KeyString) MustKey() []byte {
	return []byte(key.Value)
}

type KeyUint64 struct {
	Value uint64
}

func (key *KeyUint64) Key() ([]byte, error) {
	buf, err := GetBigEndianBytesUint64(key.Value)
	if err != nil {
		return nil, err
	}
	defer ReleaseBytesBuffer(buf)
	return buf.Bytes(), nil
}

func (key *KeyUint64) MustKey() []byte {
	buf := MustGetBigEndianBytesUint64(key.Value)
	defer ReleaseBytesBuffer(buf)
	return buf.Bytes()
}

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

type TableOptions struct {
	Enc Encoder
	Sub subspace.Subspace
}
