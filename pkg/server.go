package pkg

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

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
	hub   Hub
}

// NewServer returns a new server.
func NewServer(store *Store, homeTempl *template.Template) *Server {
	globalCh := make(chan Message)
	hub := NewHub(globalCh)

	server := &Server{
		store: *store,
		hub:   *hub,
	}

	// Start the clients hub, wait for any message.
	go hub.Start()

	// Listen on /ws to handle WebSocket messages.
	http.HandleFunc("/ws", server.wsHandler(getWsUpgrader(), hub))

	// Serve HTML on root url
	http.HandleFunc("/", server.home(homeTempl))

	return server
}

// Listen on given address
func (s *Server) Listen(addr *string) error {
	fmt.Println("Listenning on http://0.0.0.0" + *addr)
	return http.ListenAndServe(*addr, nil)
}

func (s *Server) home(homeTempl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := homeTempl.Execute(w, nil); err != nil {
			log.Println(err)
		}
	}
}

func (s *Server) wsHandler(wsupgrader *websocket.Upgrader, hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := wsupgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatalf("failed to upgrade websocket: %v\n", err)
		}

		// new client instance
		client := NewClient(conn)

		// let the client handle quit/send/recv channels
		go client.Handle()

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

	// Message doesn't include ListID
	if msg.Method == "GetLists" {
		for _, list := range s.store.AllLists() {
			c.Send(Message{
				Method: "AddList",
				Item:   *list,
			})
		}
		return
	}

	if len(msg.Item.ListID) == 0 {
		fmt.Println(msg.Method, "called without listID")
		return
	}

	// Handle message action
	switch msg.Method {
	case "GetItems":
		c.listID = msg.ListID
		fmt.Println("client", c.ID, "use list", msg.ListID)
		for _, item := range s.store.FindAll(msg.ListID) {
			s.hub.sendToListClients(msg.ListID, Message{
				Method: "AddItem",
				Item:   *item,
			})
		}

	case "UpdateItem":
		fallthrough
	case "AddItem":
		item := s.store.Create(&msg.Item)
		if item == nil {
			return
		}
		s.hub.sendToListClients(msg.ListID, Message{
			Method: "AddItem",
			Item:   *item,
		})

	case "DeleteItem":
		err := s.store.Delete(msg.ListID, msg.ItemID)
		if err != nil {
			return
		}
		s.hub.sendToListClients(msg.ListID, Message{
			Method: "DeleteItem",
			Item:   msg.Item,
		})
	}
}

func getWsUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
}
