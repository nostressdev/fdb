package tests

import (
	"fmt"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	gen_fdb "github.com/nostressdev/fdb/tests/generated"
	"log"
	"testing"
)

func Test_CreateTable(t *testing.T) {
	table, err := gen_fdb.NewUsersTable(toUsers)
	AssertError(t, err)
	row := &gen_fdb.UsersTableRow{Man: gen_fdb.User{ID: "id", Age: 57}, Ts: 123}
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
	AssertError(t, err)

	log.Println("3")

	future, err := db.ReadTransact(func(tr fdb.ReadTransaction) (interface{}, error) {
		return table.Get(tr, &gen_fdb.UsersTablePK{Ts: row.Ts, ManID: row.Man.ID})
	})
	AssertError(t, err)
	if future == nil {
		AssertError(t, fmt.Errorf("future is nil"))
	}

	log.Println("4")

	resRow, err := future.(*gen_fdb.FutureUsersTableRow).Get()
	AssertError(t, err)
	AssertEqual(t, row, resRow, "equal res")
}
