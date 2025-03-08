package nclient

import (
	"context"
	"log"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

var JobStream jetstream.Stream

var streamConfig = jetstream.StreamConfig{
	Name:       "jobs",                    // Nombre del stream
	Subjects:   []string{"jobs.*"},        // Sujetos que el stream observará
	Storage:    jetstream.FileStorage,     // Tipo de almacenamiento: archivo
	Replicas:   1,                         // Número de réplicas
	Retention:  jetstream.WorkQueuePolicy, // Política de retención
	Discard:    jetstream.DiscardNew,      // Descartar nuevos mensajes si se supera el límite
	MaxMsgs:    2000,                      // Límite de mensajes en el stream
	MaxBytes:   256 * 1024 * 1024,         // Tamaño máximo total del stream (256MB)
	MaxAge:     7 * 24 * time.Hour,        // Tiempo de vida (TTL) de los mensajes (7 días)
	MaxMsgSize: -1,                        // Sin límite en el tamaño de mensajes
}

func CreateJobStream() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	stream, err := JS.CreateOrUpdateStream(ctx, streamConfig)
	if err != nil {
		if err == jetstream.ErrStreamNameAlreadyInUse {
			log.Printf("Stream already exists!")
		} else {
			log.Fatal(err)
		}
	}
	JobStream = stream
}
