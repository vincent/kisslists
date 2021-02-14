package pkg

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Message holds any communication message.
type Message struct {
	Item
	Method string `json:"method"`
}

// Server holds the server implementation.
type Server struct {
	store Store
	hub   hub
}

// NewServer returns a new server.
func NewServer(store *Store) *Server {
	globalCh := make(chan Message)
	globalQuit := make(chan struct{})
	hub := NewHub(globalCh, globalQuit)
	defer close(globalQuit)

	server := &Server{
		store: *store,
		hub:   *hub,
	}

	// Start the clients hub, wait for any message.
	go hub.Start()

	// Listen on /ws to handle WebSocket messages.
	http.HandleFunc("/ws", server.wsHandler(getWsUpgrader(), globalQuit, hub))

	// Serve HTML on root url
	http.HandleFunc("/", ServeHome())

	return server
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
		client.OnReceive(func(msg Message) {
			go s.eventHandler(client, msg)
		})

		// register the client in the hub to
		// - listen to incoming messages
		// - send quit message
		// - handle disconnect
		hub.Register(client)

		client.send <- Message{
			Method: "Ping",
		}
	}
}

// handle message reception by a client
// use the store to answer each methods
func (s *Server) eventHandler(c *Client, msg Message) {
	// fmt.Printf("%+v : %+v\n", event, msg)

	switch msg.Method {
	case "GetItems":
		c.listID = msg.ListID
		fmt.Println("client", c.id, "use list", msg.ListID)
		for _, item := range s.store.GetItems(msg.ListID) {
			s.hub.sendToListClients(msg.ListID, Message{
				Method: "AddItem",
				Item:   *item,
			})
		}

	case "AddItem":
		item := s.store.AddItem(&msg.Item)
		if item == nil {
			return
		}
		s.hub.sendToListClients(msg.ListID, Message{
			Method: "AddItem",
			Item:   *item,
		})
	}
}

func (s *Server) Listen(addr *string) error {
	fmt.Println("Listenning on http://localhost:" + *addr)
	return http.ListenAndServe(*addr, nil)
}

func getWsUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
}
