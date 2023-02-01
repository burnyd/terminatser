package natsresponse

import (
	"log"

	"github.com/nats-io/nats.go"
)

func NatsResponse(natsurl, subject string) {
	nc, err := nats.Connect(natsurl)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	sub, _ := nc.Subscribe(subject, func(msg *nats.Msg) {
		msg.Respond([]byte("hello, I have responded to you "))
	})

	//Add a goroutine here somewhere whoever calls this thing?
	// Possible map string interface for the response?

	sub.Unsubscribe()
}
