package tests

import (
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	generated "github.com/nostressdev/fdb/gen/tests/generated"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CreateTable(t *testing.T) {
	table, err := generated.NewUsersTable(toUsers)
	assert.Nil(t, err)
	row := &generated.UsersTableRow{Man: generated.User{ID: "id", Age: true}, Ts: 123}
	fdb.MustAPIVersion(600)
	db := fdb.MustOpenDefault()

	_, err = db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		err := table.Insert(tr, row)
		return nil, err
	})
	assert.Nil(t, err)

	future, err := db.ReadTransact(func(tr fdb.ReadTransaction) (interface{}, error) {
		return table.Get(tr, &generated.UsersTablePK{Ts: row.Ts, ManID: row.Man.ID})
	})
	assert.Nil(t, err)
	assert.NotNil(t, future, "future is not nil")

	resRow, err := future.(*generated.FutureUsersTableRow).Get()
	assert.Nil(t, err)
	assert.Equal(t, row, resRow, "equal res")

	_, err = db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		return nil, table.Delete(tr, &generated.UsersTablePK{Ts: row.Ts, ManID: row.Man.ID})
	})
	assert.Nil(t, err)

	future, err = db.ReadTransact(func(tr fdb.ReadTransaction) (interface{}, error) {
		return table.Get(tr, &generated.UsersTablePK{Ts: row.Ts, ManID: row.Man.ID})
	})
	assert.Nil(t, err)

	resRow, err = future.(*generated.FutureUsersTableRow).Get()
	assert.Nil(t, err)
	assert.Nil(t, resRow)
}
