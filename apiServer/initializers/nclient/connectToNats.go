package nclient

import (
	"faas/models"
	"os"

	"github.com/nats-io/nats.go"
)

var Client models.NATSClient

func ConnectToNats() {
	// connect to nats server
	natsURL := os.Getenv("NATS_URL")
	var nc *nats.Conn
	var err error
	if natsURL == "" {
		nc, err = nats.Connect(nats.DefaultURL)
	} else {
		nc, err = nats.Connect(natsURL)
	}
	if err != nil {
		panic("Failed to connect Nats")
	}

	Client.Conn = nc
}
