package tests

import (
	"testing"

	"github.com/matryer/is"
	"github.com/vincent/kisslists/pkg"
)

func TestSqliteStore_Bootstrap(t *testing.T) {
	// Arrange
	is := is.New(t)
	testdb := NewTestDB()
	defer testdb.Dispose()

	// Act
	store := pkg.NewStore(testdb.db)
	store.Bootstrap()

	// Assert
	_, err := testdb.db.Prepare(`CREATE TABLE ListItems (id INTEGER)`)
	is.True(err != nil)
}

func TestSqliteStore_Create(t *testing.T) {
	// Arrange
	is := is.New(t)
	testdb := NewTestDB()
	defer testdb.Dispose()

	store := pkg.NewStore(testdb.db)
	store.Bootstrap()

	// Act
	item := store.Create(&pkg.Item{
		ListID:    "A list",
		Text:      "Some text",
		IsChecked: true,
	})

	// Assert
	is.True(item != nil)
	is.True(item.ItemID > 0)
	is.Equal(item.ListID, "A list")
	is.Equal(item.Text, "Some text")
	is.Equal(item.IsChecked, true)
}

func TestSqliteStore_Update(t *testing.T) {
	// Arrange
	is := is.New(t)
	testdb := NewTestDB()
	defer testdb.Dispose()

	store := pkg.NewStore(testdb.db)
	store.Bootstrap()

	// Act
	item := store.Create(&pkg.Item{
		ListID:    "A list",
		Text:      "Some text",
		IsChecked: true,
	})

	store.Create(&pkg.Item{
		ListID:    "A list",
		Text:      "Some text",
		IsChecked: false,
	})

	updated := store.Find(item.ListID, item.ItemID)

	// Assert
	is.True(item != nil)
	is.True(updated != nil)
	is.Equal(updated.IsChecked, false)
}

func TestSqliteStore_Find(t *testing.T) {
	// Arrange
	is := is.New(t)
	testdb := NewTestDB()
	defer testdb.Dispose()

	store := pkg.NewStore(testdb.db)
	store.Bootstrap()

	id := testdb.Insert(pkg.Item{
		ListID:    "A list",
		Text:      "Some text",
		IsChecked: true,
	})

	// Act
	found := store.Find("A list", id)

	// Assert
	is.True(found != nil)
	is.Equal(found.ItemID, id)
	is.Equal(found.ListID, "A list")
}

func TestSqliteStore_NotFind(t *testing.T) {
	// Arrange
	is := is.New(t)
	testdb := NewTestDB()
	defer testdb.Dispose()

	store := pkg.NewStore(testdb.db)
	store.Bootstrap()

	// Act
	found := store.Find("A List", 123)

	// Assert
	is.Equal(found, nil)
}

func TestSqliteStore_FindAll(t *testing.T) {
	// Arrange
	is := is.New(t)
	testdb := NewTestDB()
	defer testdb.Dispose()

	store := pkg.NewStore(testdb.db)
	store.Bootstrap()

	listID := "A new list"
	id1 := testdb.Insert(pkg.Item{
		ListID:    listID,
		Text:      "Some text",
		IsChecked: true,
	})
	id2 := testdb.Insert(pkg.Item{
		ListID:    listID,
		Text:      "Some other text",
		IsChecked: false,
	})

	// Act
	found := store.FindAll(listID)

	// Assert
	is.True(found != nil)
	is.Equal(found[0].ItemID, id1)
	is.Equal(found[1].ItemID, id2)
}

func TestSqliteStore_Delete(t *testing.T) {
	// Arrange
	is := is.New(t)
	testdb := NewTestDB()
	defer testdb.Dispose()

	store := pkg.NewStore(testdb.db)
	store.Bootstrap()

	id := testdb.Insert(pkg.Item{
		ListID:    "A list",
		Text:      "Some text",
		IsChecked: true,
	})

	// Act
	_ = store.Delete("A list", id)
	found := store.Find("A list", id)

	// Assert
	is.True(found == nil)
}

func TestSqliteStore_DeleteList(t *testing.T) {
	// Arrange
	is := is.New(t)
	testdb := NewTestDB()
	defer testdb.Dispose()

	store := pkg.NewStore(testdb.db)
	store.Bootstrap()

	testdb.Insert(pkg.Item{
		ListID:    "A list",
		Text:      "Some text",
		IsChecked: true,
	})

	// Act
	_ = store.DeleteList("A list")
	items := store.FindAll("A list")

	// Assert
	is.Equal(len(items), 0)
}

func TestSqliteStore_AllLists_Empty(t *testing.T) {
	// Arrange
	is := is.New(t)
	testdb := NewTestDB()
	defer testdb.Dispose()

	store := pkg.NewStore(testdb.db)
	store.Bootstrap()

	// Act
	lists := store.AllLists()

	// Assert
	is.Equal(len(lists), 0)
}

func TestSqliteStore_AllLists(t *testing.T) {
	// Arrange
	is := is.New(t)
	testdb := NewTestDB()
	defer testdb.Dispose()

	store := pkg.NewStore(testdb.db)
	store.Bootstrap()

	listIDs := []string{"abc", "def"}

	testdb.Insert(pkg.Item{
		ListID:    listIDs[0],
		Text:      "fst",
		IsChecked: true,
	})

	testdb.Insert(pkg.Item{
		ListID:    listIDs[1],
		Text:      "snd",
		IsChecked: false,
	})

	testdb.Insert(pkg.Item{
		ListID:    listIDs[0],
		Text:      "another one",
		IsChecked: false,
	})

	// Act
	lists := store.AllLists()

	// Assert
	is.Equal(len(lists), 2)
	for _, item := range lists {
		found := false
		for _, existingListID := range listIDs {
			if existingListID == item.ListID {
				found = true
				break
			}
		}
		is.True(found)
	}
}
