// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"database/sql"
	"flag"
	"io/ioutil"
	"log"
	"text/template"

	"github.com/vincent/kisslists/pkg"
)

var (
	dbfile   = flag.String("database", "./kisslists.sqlite", "SQLite database file")
	addr     = flag.String("port", ":80", "HTTP service address")
	filename string
)

func main() {
	flag.Parse()

	db, err := sql.Open("sqlite3", *dbfile)
	if err != nil {
		panic(err)
	}

	var html, _ = ioutil.ReadFile("static/frontend.html")
	var homeTpl = template.Must(template.New("").Parse(string(html)))

	store := pkg.NewStore(db)
	store.Bootstrap()

	server := pkg.NewServer(&store, homeTpl)
	if err := server.Listen(addr); err != nil {
		log.Fatal(err)
	}
}
