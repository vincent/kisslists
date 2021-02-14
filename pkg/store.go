package pkg

import (
	"context"
	"database/sql"
	"log"
	"strconv"

	// sqlite
	"github.com/georgysavva/scany/sqlscan"
	_ "github.com/mattn/go-sqlite3"
)

// Item is a stored list item
type Item struct {
	ItemID    int64  `db:"itemId"      json:"itemId"`
	ListID    string `db:"listId"      json:"listId"`
	Text      string `db:"contentText" json:"text"`
	IsChecked bool   `db:"isChecked"   json:"isChecked"`
}

// Sore is the storage backend
type Store interface {
	Bootstrap()
	GetItem(itemID int64) *Item
	GetItems(listID string) []*Item
	AddItem(item *Item) *Item
	updateItem(item *Item) *Item
}

// SqliteStore is the SQLite storge implementaion
type SqliteStore struct {
	DB *sql.DB
}

// NewStore return a new store.
// call .Bootstrap() to create the necessary tables.
func NewStore(dbfile string) Store {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		panic(err)
	}
	return &SqliteStore{
		DB: db,
	}
}

// Bootstrap the store, to create the necessary tables.
func (store *SqliteStore) Bootstrap() {
	stmt, err := store.DB.Prepare(`
		CREATE TABLE IF NOT EXISTS ListItems (
			itemId      INTEGER      PRIMARY KEY AUTOINCREMENT,
			listId      VARCHAR(128) NOT NULL,
			isChecked   TINYINT 	 DEFAULT 0,
			contentText TEXT
		) ;
		CREATE INDEX idx_list_item ON ListItems (listId, itemId);
		`)
	defer stmt.Close()
	if _, err = stmt.Exec(); err != nil {
		log.Println("created a new database")
	}
}

// GetItem returns the item matching the given ID
func (store *SqliteStore) GetItem(itemID int64) *Item {
	var items []*Item
	ctx := context.Background()
	err := sqlscan.Select(ctx, store.DB, &items,
		`SELECT itemId, listId, isChecked, contentText FROM ListItems WHERE itemId = ?`, itemID)
	if err != nil || len(items) != 1 {
		log.Println("item not found:", itemID)
		return nil
	}
	return items[0]
}

// GetItems returns all items matching the given list ID
func (store *SqliteStore) GetItems(listID string) []*Item {
	var items []*Item
	ctx := context.Background()
	err := sqlscan.Select(ctx, store.DB, &items,
		`SELECT itemId, listId, isChecked, contentText FROM ListItems WHERE listId = ?`, listID)
	if err != nil {
		log.Println(err)
	}
	return items
}

// AddItem insert (or update) the given item.
func (store *SqliteStore) AddItem(item *Item) *Item {
	listID, _ := strconv.Atoi(item.ListID)
	exists := store.getItemByText(int64(listID), item.Text)
	if exists != nil {
		exists.IsChecked = item.IsChecked
		exists = store.updateItem(exists)
		return exists
	}
	stmt, err := store.DB.Prepare(`
		INSERT INTO 
			ListItems (listId, isChecked, contentText)
			   VALUES (?, ?, ?)`)
	defer stmt.Close()
	var res sql.Result
	if res, err = stmt.Exec(item.ListID, item.IsChecked, item.Text); err != nil {
		log.Println(err)
	}
	id, err := res.LastInsertId()
	return store.GetItem(id)
}

func (store *SqliteStore) updateItem(item *Item) *Item {
	stmt, err := store.DB.Prepare(`
		UPDATE ListItems SET isChecked = ? WHERE itemId = ?`)
	defer stmt.Close()
	if _, err = stmt.Exec(item.IsChecked, item.ItemID); err != nil {
		log.Println(err)
	}
	return store.GetItem(item.ItemID)
}

func (store *SqliteStore) getItemByText(listId int64, contentText string) *Item {
	var items []*Item
	ctx := context.Background()
	err := sqlscan.Select(ctx, store.DB, &items,
		`SELECT itemId, listId, isChecked, contentText FROM ListItems WHERE contentText = ?`, contentText)
	if err != nil || len(items) != 1 {
		return nil
	}
	return items[0]
}
