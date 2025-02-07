package tests

import (
	"database/sql"
	"gomin/src/framework"
	"testing"
)

type TestDB struct {
	DB *sql.DB
	Tx *sql.Tx
}

func NewTestDB(t *testing.T) *TestDB {
	db := framework.NewDatabaseConnection()
	tx, err := db.Begin()
	if err != nil {
		panic(err.Error())
	}

	t.Cleanup(func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			t.Errorf("failed to rollback transaction: %v", err)
		}
	})

	return &TestDB{DB: db, Tx: tx}
}
