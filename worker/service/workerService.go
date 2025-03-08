package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	DefaultConsumerName = "worker"
	jobsStream          = "jobs"
	maxDeliveries       = 4
)

var executionCount int

type Config struct {
	NatsURL      string
	ConsumerName string
	WorkerName   string
}

type ExecuteResponse struct {
	ExecutionTime string `json:"executionTime"`
	Result        string `json:"result"`
	Error         string `json:"error"`
}

type Message struct {
	UserID  string `json:"userId"`
	ReplyTo string `json:"replyTo"`
	Image   string `json:"image"`
	Params  string `json:"params"`
}

type Worker struct {
	js       jetstream.JetStream
	nc       *nats.Conn
	consumer jetstream.Consumer
	config   Config
	logger   *log.Logger
}

func NewWorker(cfg Config) (*Worker, error) {
	nc, err := connectNATS(cfg.NatsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	logger := log.New(os.Stdout, fmt.Sprintf("[%s]", cfg.WorkerName), log.LstdFlags)

	return &Worker{
		js:     js,
		nc:     nc,
		config: cfg,
		logger: logger,
	}, nil
}

// MaxDeliver:  maxDeliveries,
// 		BackOff: []time.Duration{
// 			5 * time.Second,
// 			10 * time.Second,
// 			15 * time.Second,
// 		},

func (w *Worker) Setup(ctx context.Context) error {
	consumer, err := w.js.CreateOrUpdateConsumer(ctx, jobsStream, jetstream.ConsumerConfig{
		Name:        w.config.ConsumerName,
		Durable:     w.config.ConsumerName,
		Description: "Docker execution worker",

		FilterSubject: "jobs.*",
	})
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	w.consumer = consumer
	w.logger.Printf("Worker setup completed!")
	return nil
}

func (w *Worker) Start(ctx context.Context) error {
	cctx, err := w.consumer.Consume(w.handleMessage)
	if err != nil {
		return fmt.Errorf("failed to start consumer: %w", err)
	}
	defer cctx.Stop()

	<-ctx.Done()
	return nil
}

func (w *Worker) handleMessage(msg jetstream.Msg) {
	var msgData Message
	if err := json.Unmarshal(msg.Data(), &msgData); err != nil {
		w.logger.Printf("Error unmarshaling message: %v", err)
		msg.Nak()
		return
	}

	res, err := w.executeDocker(msgData.Image, msgData.Params)
	if err != nil {
		w.logger.Printf("Error executing docker: %v", err)
		msg.Nak()
		return
	}

	executionCount++

	w.logger.Printf("execution count: %d. result: %s", executionCount, string(res))

	if err := w.nc.Publish(msgData.ReplyTo, res); err != nil {
		w.logger.Printf("Error publishing response: %v", err)
		msg.Nak()
		return
	}
	msg.Ack()

}

func (w *Worker) executeDocker(image, params string) ([]byte, error) {
	start := time.Now()
	cmd := exec.Command("docker", "run", "--rm", image, params)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("docker execution failed: %w", err)
	}

	response := ExecuteResponse{
		ExecutionTime: time.Since(start).String(),
		Result:        out.String(),
	}

	return json.Marshal(response)
}

func (w *Worker) Close() {
	if w.nc != nil {
		w.nc.Close()
	}
}

func connectNATS(url string) (*nats.Conn, error) {
	if url == "" {
		url = nats.DefaultURL
	}
	return nats.Connect(url)
}
