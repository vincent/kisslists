package pkg

import (
	"log"

	"github.com/gorilla/websocket"
)

var nextID int

// Client stores a connected client.
type Client struct {
	id        int
	listID    string
	conn      *websocket.Conn
	send      chan Message
	recv      chan Message
	quit      chan struct{}
	onReceive func(msg Message)
}

// NewClient creates a new connected client.
func NewClient(conn *websocket.Conn) *Client {
	nextID++
	return &Client{
		id:        nextID,
		conn:      conn,
		quit:      make(chan struct{}),
		send:      make(chan Message),
		recv:      make(chan Message),
		onReceive: func(msg Message) {},
	}
}

// Close the client connection.
func (c *Client) Close() {
	close(c.quit)
}

// OnReceive define the receiving callback.
func (c *Client) OnReceive(f func(msg Message)) {
	c.onReceive = f
}

func (c *Client) handle() {
	for {
		select {
		case <-c.quit:
			if err := c.conn.Close(); err != nil {
				log.Printf("client %d connection close error %v\n", c.id, err)
			}
			return
		case n := <-c.send:
			if err := c.conn.WriteJSON(n); err != nil {
				log.Println("ws write error:", err)
				return
			}
		case m := <-c.recv:
			go c.onReceive(m)
		}
	}
}
