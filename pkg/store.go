package pkg

import (
	"database/sql"
	"errors"
	"log"

	// sqlite
	_ "github.com/mattn/go-sqlite3"
)

// Item is a stored list item
type Item struct {
	ItemID    int64  `db:"itemId"      json:"itemId"`
	ListID    string `db:"listId"      json:"listId"`
	Text      string `db:"contentText" json:"text"`
	IsChecked bool   `db:"isChecked"   json:"isChecked"`
}

// Store is the storage backend
type Store interface {
	Bootstrap()
	Find(listID string, itemID int64) *Item
	FindAll(listID string) []*Item
	Create(item *Item) *Item
	update(item *Item) *Item
	Delete(listID string, itemID int64) error
	AllLists() []*Item
}

// SqliteStore is the SQLite storge implementaion
type SqliteStore struct {
	DB *sql.DB
}

// NewStore return a new store.
// call .Bootstrap() to create the necessary tables.
func NewStore(db *sql.DB) Store {
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

// Find returns the item matching the given ID
func (store *SqliteStore) Find(listID string, itemID int64) *Item {
	rows, _ := store.DB.Query(
		`SELECT itemId, listId, isChecked, contentText FROM ListItems WHERE listId = ? AND itemId = ?`, listID, itemID)
	items, err := store.selectItems(rows)
	defer rows.Close()

	if err != nil || len(items) != 1 {
		return nil
	}
	return items[0]
}

// FindAll returns all items matching the given list ID
func (store *SqliteStore) FindAll(listID string) []*Item {
	rows, _ := store.DB.Query(
		`SELECT itemId, listId, isChecked, contentText FROM ListItems WHERE listId = ?`, listID)
	items, _ := store.selectItems(rows)
	defer rows.Close()
	return items
}

// Create insert (or update) the given item.
func (store *SqliteStore) Create(item *Item) *Item {
	var exists *Item
	if item.ItemID == 0 && len(item.Text) == 0 {
		return nil
	}
	if item.ItemID > 0 {
		exists = store.Find(item.ListID, item.ItemID)
	} else {
		exists = store.getItemByText(item.ListID, item.Text)
	}
	if exists != nil {
		exists.IsChecked = item.IsChecked
		exists = store.update(exists)
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
	return store.Find(item.ListID, id)
}

// Delete the given item
func (store *SqliteStore) Delete(listID string, itemID int64) error {
	exists := store.Find(listID, itemID)
	if exists == nil {
		return errors.New("No such item")
	}
	stmt, err := store.DB.Prepare(`
		DELETE FROM ListItems WHERE itemId = ?`)
	defer stmt.Close()
	if _, err = stmt.Exec(itemID); err != nil {
		log.Println(err)
	}
	return err
}

func (store *SqliteStore) update(item *Item) *Item {
	stmt, err := store.DB.Prepare(`
		UPDATE ListItems SET isChecked = ? WHERE itemId = ?`)
	defer stmt.Close()
	if _, err = stmt.Exec(item.IsChecked, item.ItemID); err != nil {
		log.Println(err)
	}
	return store.Find(item.ListID, item.ItemID)
}

func (store *SqliteStore) getItemByText(listID string, contentText string) *Item {
	rows, _ := store.DB.Query(
		`SELECT itemId, listId, isChecked, contentText FROM ListItems WHERE listId = ? AND contentText = ?`, listID, contentText)
	items, err := store.selectItems(rows)
	defer rows.Close()
	if err != nil || len(items) != 1 {
		return nil
	}
	return items[0]
}

func (store *SqliteStore) selectItems(rows *sql.Rows) ([]*Item, error) {
	result := []*Item{}

	var itemID int64
	var listID string
	var text string
	var isChecked bool

	for rows.Next() {
		err := rows.Scan(&itemID, &listID, &isChecked, &text)
		if err != nil {
			return result, err
		}

		result = append(result, &Item{
			ItemID:    itemID,
			ListID:    listID,
			Text:      text,
			IsChecked: isChecked,
		})
	}

	return result, nil
}

func (store *SqliteStore) AllLists() []*Item {
	rows, _ := store.DB.Query("SELECT DISTINCT listId FROM ListItems")
	var lists []*Item
	var listID string
	for rows.Next() {
		err := rows.Scan(&listID)
		if err != nil {
			return nil
		}
		lists = append(lists, &Item{ListID: listID})
	}

	return lists
}
