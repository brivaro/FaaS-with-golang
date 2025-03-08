package nclient

import (
	"context"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

func CreateFunctionKV() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	kvFunctions, err := JS.CreateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket: "functions",
	})
	if err != nil {
		panic("Error creating users Key value store")
	}
	Client.KvFunctions = kvFunctions
}
