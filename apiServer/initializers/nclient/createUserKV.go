package nclient

import (
	"context"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

func CreateUserKV() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	kvUsers, err := JS.CreateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket: "users",
	})
	if err != nil {
		panic("Error creating users Key value store")
	}
	Client.KvUsers = kvUsers
}
