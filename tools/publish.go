package main

import (
	"log"

	nats "github.com/nats-io/go-nats"
)

const (
	topicUnknown = "out.unknown"
	testMessage  = `{	\"id\":\"5c59ee1c0e360800016fc255\",
						\"pushed\":0,
						\"device\":\"Random-Integer-Generator01\",
						\"created\":1549397532000,
						\"modified\":0,
						\"origin\":1549397532000,
						\"schedule\":null,
						\"event\":null,
						\"readings\":
							[{	\"id\":\"5c59ee1c0e360800016fc256\",
								\"pushed\":0,
								\"created\":1549397532000,
								\"origin\":1549397532000,
								\"modified\":0,
								\"device\":\"Random-Integer-Generator01\",
								\"name\":\"RandomValue_Int32\",
								\"value\":\"116177668\"
							}]
					}`
)

func main() {

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	if err := nc.Publish(topicUnknown, []byte(testMessage)); err != nil {
		log.Fatal(err)
	}
	// Make sure the message goes through before we close
	nc.Flush()
}
