package pkg

import (
	"log"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients map[int]*Client
	// broadcast channel on which the hub will receiver messages
	// and broadcast them all clients
	broadcast chan Message
	// register channel on which the hub will register clients
	register chan *Client
	// deregister channel on which the hub will deregister clients
	deregister chan *Client
	// quit channel
	quit chan int
}

// NewHub return a new clients hub
func NewHub(ch chan Message) *Hub {
	return &Hub{
		clients:    make(map[int]*Client),
		register:   make(chan *Client),
		deregister: make(chan *Client),
		quit:       make(chan int),
		broadcast:  ch,
	}
}

func (h *Hub) Start() {
	for {
		select {
		case message := <-h.broadcast:
			for _, client := range h.clients {
				client.Send(message)
			}
		case client := <-h.register:
			h.push(client)
			go h.watchDisconnect(client)
		case client := <-h.deregister:
			h.delete(client)
		case _ = <-h.quit:
			return
		}
	}
}

// read each message received on the client connection
// and handle disconnect or route to the client recv channel
func (h *Hub) watchDisconnect(client *Client) {
	for {
		var msg Message
		err := client.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsCloseError(err) {
				log.Println("client disconnected", client.ID)
			}
			client.Close()
			h.Deregister(client)
			return
		}
		client.Receive(msg)
	}
}

func (h *Hub) Stop() {
	close(h.quit)
}

func (h *Hub) Register(c *Client) {
	log.Printf("client %d connected.\n", c.ID)
	h.register <- c
}

func (h *Hub) Deregister(c *Client) {
	log.Printf("client %d disconnected.\n", c.ID)
	h.deregister <- c
}

func (h *Hub) Iter(f func(*Client)) {
	for _, client := range h.clients {
		f(client)
	}
}

func (h *Hub) push(client *Client) {
	h.clients[client.ID] = client
}

func (h *Hub) delete(client *Client) {
	delete(h.clients, client.ID)
}

func (h *Hub) sendToListClients(listID string, m Message) {
	for _, client := range h.clients {
		if client.listID == listID {
			client.Send(m)
		}
	}
}
