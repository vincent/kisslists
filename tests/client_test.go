package tests

import (
	"sync"
	"testing"

	"github.com/matryer/is"
	"github.com/vincent/sharedlists/pkg"
)

func TestClient_NewClient(t *testing.T) {
	// Arrange
	is := is.New(t)

	// Act
	c1 := pkg.NewClient(nil)
	c2 := pkg.NewClient(nil)
	c3 := pkg.NewClient(nil)

	// Assert
	is.True(c1.ID != c2.ID)
	is.True(c2.ID != c3.ID)
}

func TestClient_OnReceive(t *testing.T) {
	// Arrange
	is := is.New(t)
	c := pkg.NewClient(nil)

	var wg sync.WaitGroup
	var rcvHasBeenCalled bool
	var msgReceived *pkg.Message

	// Act
	wg.Add(1)
	c.OnReceive(func(msg pkg.Message) {
		rcvHasBeenCalled = true
		msgReceived = &msg
		wg.Done()
	})
	go c.Handle()
	c.Receive(pkg.Message{Method: "TEST"})
	wg.Wait()

	// Assert
	is.True(rcvHasBeenCalled)
	is.True(msgReceived != nil)
	is.Equal(msgReceived.Method, "TEST")
}
