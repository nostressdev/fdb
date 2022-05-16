package generate

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
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
//		Fields:
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
//		PK: ["ID", "Points"]
// }

type UsersTable struct {
	DB       fdb.Database
	Subspace subspace.Subspace
}

type User struct {
	ID     string `json:"id"`
	Points uint64 `json:"points"`
}

func NewUser(jsonBytes []byte) (*User, error) {
	model := &User{}
	if err := json.Unmarshal(jsonBytes, model); err != nil {
		return nil, err
	}
	return model, nil
}

func MustNewUser(jsonBytes []byte) *User {
	model := &User{}
	if err := json.Unmarshal(jsonBytes, model); err != nil {
		panic(err)
	}
	return model
}

func (model *User) ToJson() ([]byte, error) {
	return json.Marshal(model)
}

func (model *User) MustToJson() []byte {
	if value, err := json.Marshal(model); err != nil {
		panic(err)
	} else {
		return value
	}
}

type UsersTablePK struct {
	ID     Key // string
	Points Key // uint64
}

type UsersTablePKBytes struct {
	IDBytes     []byte // string
	PointsBytes []byte // uint64
}

func (table *UsersTable) Get(tr fdb.ReadTransaction, pk *UsersTablePK) (*User, error) {
	pkBytes := &UsersTablePKBytes{}
	if pkIDBytes, err := pk.ID.Key(); err != nil {
		return nil, err
	} else {
		pkBytes.IDBytes = pkIDBytes
	}
	if pkPointsBytes, err := pk.Points.Key(); err != nil {
		return nil, err
	} else {
		pkBytes.PointsBytes = pkPointsBytes
	}

	valueJson, err := tr.Get(table.Subspace.Sub(pkBytes.IDBytes, pkBytes.PointsBytes)).Get()
	if err != nil {
		return nil, err
	}
	return NewUser(valueJson)
}

func (table *UsersTable) MustGet(tr fdb.ReadTransaction, pk *UsersTablePK) *User {
	pkBytes := &UsersTablePKBytes{}
	pkBytes.IDBytes = pk.ID.MustKey()
	pkBytes.PointsBytes = pk.Points.MustKey()

	return MustNewUser(tr.Get(table.Subspace.Sub(pkBytes.IDBytes, pkBytes.PointsBytes)).MustGet())
}

func (table *UsersTable) Insert(tr fdb.Transaction, pk *UsersTablePK, model *User) error {
	pkBytes := &UsersTablePKBytes{}
	if pkIDBytes, err := pk.ID.Key(); err != nil {
		return err
	} else {
		pkBytes.IDBytes = pkIDBytes
	}
	if pkPointsBytes, err := pk.Points.Key(); err != nil {
		return err
	} else {
		pkBytes.PointsBytes = pkPointsBytes
	}

	valueJson, err := model.ToJson()
	if err != nil {
		return err
	}
	tr.Set(table.Subspace.Sub(pkBytes.IDBytes, pkBytes.PointsBytes), valueJson)
	return nil
}

func (table *UsersTable) MustInsert(tr fdb.Transaction, pk *UsersTablePK, model *User) {
	pkBytes := &UsersTablePKBytes{}
	pkBytes.IDBytes = pk.ID.MustKey()
	pkBytes.PointsBytes = pk.Points.MustKey()

	valueJson := model.MustToJson()
	tr.Set(table.Subspace.Sub(pkBytes.IDBytes, pkBytes.PointsBytes), valueJson)
}

func (table *UsersTable) Delete(tr fdb.Transaction, pk *UsersTablePK) error {
	pkBytes := &UsersTablePKBytes{}
	if pkIDBytes, err := pk.ID.Key(); err != nil {
		return err
	} else {
		pkBytes.IDBytes = pkIDBytes
	}
	if pkPointsBytes, err := pk.Points.Key(); err != nil {
		return err
	} else {
		pkBytes.PointsBytes = pkPointsBytes
	}

	tr.Clear(table.Subspace.Sub(pkBytes.IDBytes, pkBytes.PointsBytes))
	return nil
}

func (table *UsersTable) MustDelete(tr fdb.Transaction, pk *UsersTablePK) {
	pkBytes := &UsersTablePKBytes{}
	pkBytes.IDBytes = pk.ID.MustKey()
	pkBytes.PointsBytes = pk.Points.MustKey()

	tr.Clear(table.Subspace.Sub(pkBytes.IDBytes, pkBytes.PointsBytes))
}
