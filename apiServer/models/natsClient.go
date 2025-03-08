package models

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NATSClient struct {
	Conn        *nats.Conn
	KvUsers     jetstream.KeyValue
	KvFunctions jetstream.KeyValue
}
