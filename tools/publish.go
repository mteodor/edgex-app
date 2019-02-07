package main

import (
	"log"

	nats "github.com/nats-io/go-nats"
)

func main() {

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	if err := nc.Publish("out.unknown", []byte("All is Well")); err != nil {
		log.Fatal(err)
	}
	// Make sure the message goes through before we close
	nc.Flush()
}
