package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/matryer/is"
	"github.com/vincent/sharedlists/pkg"
)

func TestHub_Iter(t *testing.T) {
	// Arrange
	is := is.New(t)
	hub := pkg.NewHub(nil)
	calls := 0

	// Act
	go hub.Start()
	hub.Register(pkg.NewClient(NewTestWebsocket(t, echoHandler())))
	hub.Register(pkg.NewClient(NewTestWebsocket(t, echoHandler())))
	time.Sleep(10 * time.Millisecond) // thats bad
	hub.Iter(func(c *pkg.Client) {
		calls++
	})

	// Assert
	is.Equal(calls, 2)
}

/******************************/

func NewTestWebsocket(t *testing.T, handler http.HandlerFunc) *websocket.Conn {
	// Create test server with the echo handler.
	s := httptest.NewServer(handler)

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}

	return ws
}

var upgrader = websocket.Upgrader{}

func echoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				break
			}
			err = c.WriteMessage(mt, message)
			if err != nil {
				break
			}
		}
	}
}
