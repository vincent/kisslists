package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/vincent/kisslists/embed"
	"github.com/vincent/kisslists/pkg"
)

var (
	dbfile   = flag.String("database", "/tmp/kisslists.sqlite", "SQLite database file")
	addr     = flag.String("port", ":80", "HTTP service address")
	filename string
)

func main() {
	flag.Parse()

	// Test if the database is writable
	if touch(*dbfile) != nil {
		panic(fmt.Errorf("%v is not usable", *dbfile))
	}

	// Open database file
	db, err := sql.Open("sqlite3", *dbfile)
	if err != nil {
		panic(err)
	}

	// Load HTML from generated assets
	html := string(embed.Get("/frontend.html"))

	// Let the template include other assets with {{include}}
	var funcMap = map[string]interface{}{"include": embed.Include}
	homeTpl := template.Must(template.New("").Funcs(funcMap).Parse(string(html)))

	// Initialise database, create it if needed
	store := pkg.NewStore(db)
	store.Bootstrap()

	// Server frontend & websockets
	server := pkg.NewServer(&store, homeTpl)
	if err := server.Listen(addr); err != nil {
		log.Fatal(err)
	}
}

// touch as in unix
func touch(name string) error {
	file, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, 0644)
	defer file.Close()
	if err != nil {
		return err
	}
	return err
}
