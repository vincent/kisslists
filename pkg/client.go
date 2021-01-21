package pkg

import (
	"log"

	"github.com/gorilla/websocket"
)

var NextID int

type client struct {
	id        int
	listID    string
	conn      *websocket.Conn
	send      chan message
	recv      chan message
	quit      chan struct{}
	onReceive func(msg message)
}

func NewClient(conn *websocket.Conn) *client {
	NextID++
	return &client{
		id:        NextID,
		conn:      conn,
		quit:      make(chan struct{}),
		send:      make(chan message),
		recv:      make(chan message),
		onReceive: func(msg message) {},
	}
}

func (c *client) Close() {
	close(c.quit)
}

func (c *client) OnReceive(f func(msg message)) {
	c.onReceive = f
}

func (c *client) handle() {
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

// func (c *client) listen() {
// 	for {
// 		var msg *message
// 		err := c.conn.ReadJSON(&msg)
// 		if err != nil {
// 			fmt.Println("read:", err)
// 			continue
// 		}
// 		c.recv <- *msg
// 	}
// }
