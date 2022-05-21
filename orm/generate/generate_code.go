package generate

import (
	"fmt"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
	"github.com/nostressdev/fdb/lib"
)

// код в библе:

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
	Enc lib.Encoder
	Sub subspace.Subspace
}

func NewUsersTable(opts ...lib.TableOptions) (*UsersTable, error) {
	table := &UsersTable{}
	for _, opt := range opts {
		if opt.Enc != nil {
			table.Enc = opt.Enc
		}
		if opt.Sub != nil {
			table.Sub = opt.Sub
		}
	}
	if table.Enc == nil {
		return nil, fmt.Errorf("encoder is nil")
	}
	if table.Sub == nil {
		return nil, fmt.Errorf("subspace is nil")
	}
	return table, nil
}

type User struct {
	ID     string `json:"id"`
	Points uint64 `json:"points"`
}

type UsersTableRow struct {
	User User   `json:"user"`
	Ts   uint64 `json:"ts"`
}

type FutureUsersTableRow struct {
	Enc    lib.Encoder
	Future fdb.FutureByteSlice
}

func (future *FutureUsersTableRow) NewUsersTableRow(value []byte) (*UsersTableRow, error) {
	if row, err := future.Enc.Decode(value); err != nil {
		return nil, err
	} else {
		return row.(*UsersTableRow), nil
	}
}

func (future *FutureUsersTableRow) Get() (*UsersTableRow, error) {
	value, err := future.Future.Get()
	if err != nil {
		return nil, err
	}
	return future.NewUsersTableRow(value)
}

func (future *FutureUsersTableRow) MustGet() *UsersTableRow {
	value, err := future.Get()
	if err != nil {
		panic(err)
	}
	return value
}

type UsersTablePK struct {
	ID lib.Key // string
	Ts lib.Key // uint64
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

func (table *UsersTable) Get(tr fdb.ReadTransaction, pk *UsersTablePK) (*FutureUsersTableRow, error) {
	key, err := pk.Pack()
	if err != nil {
		return nil, err
	}

	future := tr.Get(table.Sub.Sub(key))
	return &FutureUsersTableRow{Future: future, Enc: table.Enc}, nil
}

func (table *UsersTable) MustGet(tr fdb.ReadTransaction, pk *UsersTablePK) *FutureUsersTableRow {
	future, err := table.Get(tr, pk)
	if err != nil {
		panic(err)
	}
	return future
}

func (table *UsersTable) Insert(tr fdb.Transaction, model *UsersTableRow) error {
	pk := &UsersTablePK{
		ID: &lib.KeyString{Value: model.User.ID},
		Ts: &lib.KeyUint64{Value: model.Ts},
	}
	key, err := pk.Pack()
	if err != nil {
		return err
	}

	value, err := table.Enc.Encode(model)
	if err != nil {
		return err
	}
	tr.Set(table.Sub.Sub(key), value)
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

	tr.Clear(table.Sub.Sub(key))
	return nil
}

func (table *UsersTable) MustDelete(tr fdb.Transaction, pk *UsersTablePK) {
	err := table.Delete(tr, pk)
	if err != nil {
		panic(err)
	}
}
