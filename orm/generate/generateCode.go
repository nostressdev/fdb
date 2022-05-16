package generate

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
	"github.com/nostressdev/fdb/orm/scheme"
	"sync"
)

func generateCode(config *scheme.GeneratorConfig) error {
	panic("implement me")
}

// код в библе:
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
	value string
}

func (key *KeyString) Key() ([]byte, error) {
	return []byte(key.value), nil
}

func (key *KeyString) MustKey() []byte {
	return []byte(key.value)
}

type KeyUint64 struct {
	value uint64
}

func (key *KeyUint64) Key() ([]byte, error) {
	buf, err := GetBigEndianBytesUint64(key.value)
	if err != nil {
		return nil, err
	}
	defer ReleaseBytesBuffer(buf)
	return buf.Bytes(), nil
}

func (key *KeyUint64) MustKey() []byte {
	buf := MustGetBigEndianBytesUint64(key.value)
	defer ReleaseBytesBuffer(buf)
	return buf.Bytes()
}

// Ниже мое представление о том, как будет выглядеть сгенерированный код
// Это необходимо мне для понимания того что писать)
// схема GeneratorConfig
// scheme.Model
// {
//		name: "User",
//		Fields: // такой же вопрос как и к колонкам
//		[
//			{
//				name: "ID",
//				type: "string"
//				DefaultValue: ""
//			},
//			{
//				name: "Points",
//				type: "uint64"
//				DefaultValue: "0"
//			}
//		]
// }
// scheme.Table
// {
//		name: "Users",
//      columns: ["user": {"User"}, "ts": {"uint64"}] // почему Columns это мэп, а не список, если имя к колонки есть в структуре
//		PK: ["ts", "user.ID"]
// }

type UsersTable struct {
	DB       fdb.Database
	Subspace subspace.Subspace
}

type User struct {
	ID     string `json:"id"`
	Points uint64 `json:"points"`
}

type UsersTableRow struct {
	User User   `json:"user"`
	Ts   uint64 `json:"ts"`
}

func NewUsersTableRow(jsonBytes []byte) (*UsersTableRow, error) {
	model := &UsersTableRow{}
	if err := json.Unmarshal(jsonBytes, model); err != nil {
		return nil, err
	}
	return model, nil
}

func MustNewUsersTableRow(jsonBytes []byte) *UsersTableRow {
	model := &UsersTableRow{}
	if err := json.Unmarshal(jsonBytes, model); err != nil {
		panic(err)
	}
	return model
}

type FutureUserTableRow struct {
	Future fdb.FutureByteSlice
}

func (future *FutureUserTableRow) Get() (*UsersTableRow, error) {
	valueJson, err := future.Future.Get()
	if err != nil {
		return nil, err
	}
	return NewUsersTableRow(valueJson)
}

func (future *FutureUserTableRow) MustGet() *UsersTableRow {
	return MustNewUsersTableRow(future.Future.MustGet())
}

func (model *UsersTableRow) ToJson() ([]byte, error) {
	return json.Marshal(model)
}

func (model *UsersTableRow) MustToJson() []byte {
	if value, err := json.Marshal(model); err != nil {
		panic(err)
	} else {
		return value
	}
}

type UsersTablePK struct {
	ID Key // string
	Ts Key // uint64
}

func (pk *UsersTablePK) Pack() ([]tuple.TupleElement, error) {
	pkIDBytes, err := pk.ID.Key()
	if err != nil {
		return nil, err
	}
	pkTsBytes, err := pk.Ts.Key()
	if err != nil {
		return nil, err
	}
	return []tuple.TupleElement{pkIDBytes, pkTsBytes}, nil
}

func (table *UsersTable) Get(tr fdb.ReadTransaction, pk *UsersTablePK) (*FutureUserTableRow, error) {
	key, err := pk.Pack()
	if err != nil {
		return nil, err
	}

	future := tr.Get(table.Subspace.Sub(key))
	return &FutureUserTableRow{Future: future}, nil
}

func (table *UsersTable) MustGet(tr fdb.ReadTransaction, pk *UsersTablePK) *FutureUserTableRow {
	future, err := table.Get(tr, pk)
	if err != nil {
		panic(err)
	}
	return future
}

func (table *UsersTable) Insert(tr fdb.Transaction, model *UsersTableRow) error {
	pk := &UsersTablePK{
		ID: &KeyString{model.User.ID},
		Ts: &KeyUint64{model.Ts},
	}
	key, err := pk.Pack()
	if err != nil {
		return err
	}

	valueJson, err := model.ToJson()
	if err != nil {
		return err
	}
	tr.Set(table.Subspace.Sub(key), valueJson)
	return nil
}

func (table *UsersTable) MustInsert(tr fdb.Transaction, model *UsersTableRow) {
	err := table.Insert(tr, model)
	if err != nil {
		panic(err)
	}
}

func (table *UsersTable) Delete(tr fdb.Transaction, pk *UsersTablePK) error {
	key, err := pk.Pack()
	if err != nil {
		return err
	}

	tr.Clear(table.Subspace.Sub(key))
	return nil
}

func (table *UsersTable) MustDelete(tr fdb.Transaction, pk *UsersTablePK) {
	err := table.Delete(tr, pk)
	if err != nil {
		panic(err)
	}
}
