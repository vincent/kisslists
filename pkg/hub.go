package pkg

import (
	"log"

	"github.com/gorilla/websocket"
)

type hub struct {
	clients map[int]*client
	// broadcast channel on which the hub will receiver messages
	// and broadcast them all clients
	broadcast chan message
	// quit channel
	quit chan struct{}
	// register channel on which the hub will register clients
	register chan *client
	// deregister channel on which the hub will deregister clients
	deregister chan *client
}

func NewHub(ch chan message, quit chan struct{}) *hub {
	return &hub{
		clients:    make(map[int]*client),
		register:   make(chan *client),
		deregister: make(chan *client),
		quit:       quit,
		broadcast:  ch,
	}
}

func (h *hub) start() {
	for {
		select {
		case <-h.quit:
			log.Println("hub: global quit")
			return
		case message := <-h.broadcast:
			for _, client := range h.clients {
				client.send <- message
			}
		case client := <-h.register:
			h.push(client)
			go h.watchDisconnect(client)
		case client := <-h.deregister:
			h.delete(client)
		}
	}
}

// read each message received on the client connection
// and handle disconnect or route to the client recv channel
func (h *hub) watchDisconnect(client *client) {
	for {
		var msg message
		err := client.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsCloseError(err) {
				log.Println("client disconnected", client.id)
			}
			client.Close()
			h.Deregister(client)
			return
		}
		client.recv <- msg
	}
}

func (h *hub) Register(c *client) {
	log.Printf("client %d connected.\n", c.id)
	h.register <- c
}

func (h *hub) Deregister(c *client) {
	log.Printf("client %d disconnected.\n", c.id)
	h.deregister <- c
}

func (h *hub) push(client *client) {
	h.clients[client.id] = client
}

func (h *hub) delete(client *client) {
	delete(h.clients, client.id)
}

func (h *hub) Iter(f func(*client)) {
	for _, client := range h.clients {
		f(client)
	}
}

func (h *hub) sendToListClients(listID string, m message) {
	for _, client := range h.clients {
		if client.listID == listID {
			client.send <- m
		}
	}
}
