package pkg

import (
	"log"

	"github.com/gorilla/websocket"
)

type hub struct {
	clients map[int]*Client
	// broadcast channel on which the hub will receiver messages
	// and broadcast them all clients
	broadcast chan Message
	// register channel on which the hub will register clients
	register chan *Client
	// deregister channel on which the hub will deregister clients
	deregister chan *Client
}

func NewHub(ch chan Message, quit chan struct{}) *hub {
	return &hub{
		clients:    make(map[int]*Client),
		register:   make(chan *Client),
		deregister: make(chan *Client),
		broadcast:  ch,
	}
}

func (h *hub) Start() {
	for {
		select {
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
func (h *hub) watchDisconnect(client *Client) {
	for {
		var msg Message
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

func (h *hub) Register(c *Client) {
	log.Printf("client %d connected.\n", c.id)
	h.register <- c
}

func (h *hub) Deregister(c *Client) {
	log.Printf("client %d disconnected.\n", c.id)
	h.deregister <- c
}

func (h *hub) push(client *Client) {
	h.clients[client.id] = client
}

func (h *hub) delete(client *Client) {
	delete(h.clients, client.id)
}

func (h *hub) Iter(f func(*Client)) {
	for _, client := range h.clients {
		f(client)
	}
}

func (h *hub) sendToListClients(listID string, m Message) {
	for _, client := range h.clients {
		if client.listID == listID {
			client.send <- m
		}
	}
}
