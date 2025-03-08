package nclient

import (
	"log"

	"github.com/nats-io/nats.go/jetstream"
)

var JS jetstream.JetStream

func CreateJetStream() {
	// create jetstream context from nats connection
	js, err := jetstream.New(Client.Conn)
	if err != nil {
		log.Fatal(err)
	}
	JS = js
}
