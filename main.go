package main

import (
	"fmt"
	"time"

	nats "github.com/nats-io/go-nats"
)

// Connect to a server

func main() {

	// Connect to a server
	nc, _ := nats.Connect(nats.DefaultURL)
	// Simple Async Subscriber
	nc.Subscribe("out.unknown", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})

	nc, _ := nats.Connect(nats.DefaultURL)

	// Simple Publisher
	nc.Publish("foo", []byte("Hello World"))

	// Simple Async Subscriber
	nc.Subscribe("foo", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})

	// Simple Sync Subscriber
	sub, err := nc.SubscribeSync("foo")
	m, err := sub.NextMsg(timeout)

	// Channel Subscriber
	ch := make(chan *nats.Msg, 64)
	sub, err := nc.ChanSubscribe("foo", ch)
	msg := <-ch

	// Unsubscribe
	sub.Unsubscribe()

	// Drain
	sub.Drain()

	// Requests
	msg, err := nc.Request("help", []byte("help me"), 10*time.Millisecond)

	// Replies
	nc.Subscribe("help", func(m *Msg) {
		nc.Publish(m.Reply, []byte("I can help!"))
	})

}

func closeConn() {
	// Drain connection (Preferred for responders)
	// Close() not needed if this is called.
	nc.Drain()

	// Close connection
	nc.Close()
}
