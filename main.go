// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/vincent/sharedlists/pkg"
)

var (
	dbfile   = flag.String("database", "./sharedlists.sqlite", "SQLite database file")
	addr     = flag.String("port", ":80", "HTTP service address")
	filename string
)

func main() {
	flag.Parse()

	store := pkg.NewStore(*dbfile)
	store.Bootstrap()

	wsh := pkg.NewWebSocketHandler(&store)

	http.HandleFunc("/ws", wsh.ServeWebsocket)
	http.HandleFunc("/", pkg.ServeHome)
	fmt.Println("Listenning on http://localhost:" + *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
