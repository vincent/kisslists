package tests

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vincent/sharedlists/pkg"
)

type testDB struct {
	file string
	db   *sql.DB
}

func NewTestDB() *testDB {
	file := fmt.Sprintf("db-%d", rand.Int())
	os.Remove(file)
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		fmt.Println(err)
	}
	return &testDB{
		file: file,
		db:   db,
	}
}

func (t *testDB) Dispose() {
	t.db.Close()
	os.Remove(t.file)
}

func (t *testDB) Insert(item pkg.Item) int64 {
	stmt, err := t.db.Prepare(`
		INSERT INTO 
			ListItems (listId, isChecked, contentText)
			   VALUES (?, ?, ?)`)
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
	}

	res, err := stmt.Exec(item.ListID, item.IsChecked, item.Text)
	if err != nil {
		fmt.Println(err)
	}

	lid, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
	}

	return lid
}
