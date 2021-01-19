package pkg

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type message struct {
	Item
	Method string `json:"method"`
}

type Handler interface {
	ServeWebsocket(w http.ResponseWriter, r *http.Request)
}

type WebSocketHandler struct {
	Store Store
}

func NewWebSocketHandler(store *Store) *WebSocketHandler {
	return &WebSocketHandler{
		Store: *store,
	}
}

func (wsh *WebSocketHandler) eventHandler(c *websocket.Conn, event []byte) {
	var msg message
	err := json.Unmarshal(event, &msg)
	if err != nil {
		fmt.Println(msg)
		return
	}

	// fmt.Printf("%+v : %+v\n", event, msg)

	switch msg.Method {
	case "GetItems":
		for _, item := range wsh.Store.GetItems(msg.ListID) {
			wsh.outputMessage(c, &message{
				Method: "AddItem",
				Item:   *item,
			})
		}

	case "AddItem":
		item := wsh.Store.AddItem(&msg.Item)
		wsh.outputMessage(c, &message{
			Method: "AddItem",
			Item:   *item,
		})

	case "UpdateItem":
		item := wsh.Store.UpdateItem(&msg.Item)
		wsh.outputMessage(c, &message{
			Method: "AddItem",
			Item:   *item,
		})
	}
}

func (wsh *WebSocketHandler) ServeWebsocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			break
		}
		wsh.eventHandler(c, message)
	}
}

func (wsh *WebSocketHandler) outputMessage(c *websocket.Conn, msg *message) {
	result, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("write:", err)
	}
	err = c.WriteMessage(1, result)
	if err != nil {
		fmt.Println("write:", err)
	}
}
