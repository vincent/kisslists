package pkg

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type counter struct {
	v int
}

type message struct {
	Item
	Method string `json:"method"`
}

type Server struct {
	store Store
	hub   hub
}

func NewServer(addr *string, store *Store) {
	globalCh := make(chan message)
	globalQuit := make(chan struct{})
	hub := NewHub(globalCh, globalQuit)
	defer close(globalQuit)

	server := &Server{
		store: *store,
		hub:   *hub,
	}

	go hub.start()
	go updateCounterEvery(5*time.Second, globalCh)

	http.HandleFunc("/ws", server.wsHandler(getWsUpgrader(), globalQuit, hub))
	http.HandleFunc("/", ServeHome())

	server.Listen(addr)
}

func (s *Server) wsHandler(wsupgrader *websocket.Upgrader, globalQuit chan struct{}, hub *hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := wsupgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatalf("failed to upgrade websocket: %v\n", err)
		}

		// new client instance
		client := NewClient(conn)
		// let the client handle quit/send/recv channels
		go client.handle()

		// when client receive a message, call the eventHandler
		client.OnReceive(func(msg message) {
			go s.eventHandler(client, msg)
		})

		// register the client in the hub to
		// - listen to incoming messages
		// - send quit message
		// - handle disconnect
		hub.Register(client)

		client.send <- message{
			Method: "Ping",
		}
	}
}

// handle message reception by a client
// use the store to answer each methods
func (s *Server) eventHandler(c *client, msg message) {
	// fmt.Printf("%+v : %+v\n", event, msg)

	switch msg.Method {
	case "GetItems":
		c.listID = msg.ListID
		fmt.Println("client", c.id, "use list", msg.ListID)
		for _, item := range s.store.GetItems(msg.ListID) {
			s.hub.sendToListClients(msg.ListID, message{
				Method: "AddItem",
				Item:   *item,
			})
		}

	case "AddItem":
		item := s.store.AddItem(&msg.Item)
		s.hub.sendToListClients(msg.ListID, message{
			Method: "AddItem",
			Item:   *item,
		})

	case "UpdateItem":
		item := s.store.UpdateItem(&msg.Item)
		s.hub.sendToListClients(msg.ListID, message{
			Method: "AddItem",
			Item:   *item,
		})
	}
}

func (s *Server) Listen(addr *string) {
	fmt.Println("Listenning on http://localhost:" + *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}

func getWsUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
}

func updateCounterEvery(d time.Duration, globalCh chan message) {
	ticker := time.NewTicker(d)
	for {
		select {
		case <-ticker.C:
			globalCh <- message{
				Method: "Ping",
			}
		}
	}
}
