package tests

import (
	"fmt"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	generated "github.com/nostressdev/fdb/gen/tests/generated"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func Test_CreateTable(t *testing.T) {
	table, err := generated.NewUsersTable(toUsers)
	assert.Nil(t, err)
	row := &generated.UsersTableRow{Man: generated.User{ID: "id", Age: 57}, Ts: 123}
	fdb.MustAPIVersion(600)
	db := fdb.MustOpenDefault()

	log.Println("1")

	_, err = db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		log.Println("1.1")
		err := table.Insert(tr, row)
		log.Println("1.2")
		return nil, err
	})
	log.Println("2")
	assert.Nil(t, err)

	log.Println("3")

	future, err := db.ReadTransact(func(tr fdb.ReadTransaction) (interface{}, error) {
		return table.Get(tr, &generated.UsersTablePK{Ts: row.Ts, ManID: row.Man.ID})
	})
	assert.Nil(t, err)
	if future == nil {
		assert.Nil(t, fmt.Errorf("future is nil"))
	}

	log.Println("4")

	resRow, err := future.(*generated.FutureUsersTableRow).Get()
	assert.Nil(t, err)
	assert.Equal(t, row, resRow, "equal res")
}