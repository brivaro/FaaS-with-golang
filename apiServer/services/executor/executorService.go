package executor

import (
	"context"
	"encoding/json"
	"errors"
	"faas/models"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type ExecutorService struct {
	natsConn *nats.Conn
	js       jetstream.JetStream
}

func NewExecutorService(natsConn *nats.Conn, js jetstream.JetStream) *ExecutorService {
	return &ExecutorService{
		natsConn: natsConn,
		js:       js,
	}
}

func (e *ExecutorService) ExecuteFunction(ctx context.Context, user models.User, req ExecuteRequest) (ExecuteTask, error) {
	function, err := e.lookupFunction(req.FuncID)
	replyTo := e.natsConn.NewInbox()

	if err != nil {
		return ExecuteTask{}, fmt.Errorf("error looking for function: %v", err)
	}
	if function.UserID != user.Username {
		return ExecuteTask{}, ErrUnauthorized
	}

	task := ExecuteTask{
		UserID:  user.Username,
		ReplyTo: replyTo,
		Image:   function.Data,
		Params:  req.Parameter,
	}

	err = e.publishExecutionTask(ctx, task)

	if err != nil {
		var apiErr *jetstream.APIError
		// Max messages error exceeded
		if errors.As(err, &apiErr) && apiErr.ErrorCode == 10077 {
			return ExecuteTask{}, ErrMaxMessages
		} else {
			return ExecuteTask{}, err
		}
	}

	return task, nil
}

func (e *ExecutorService) GetResult(ctx context.Context, task ExecuteTask) (ExecuteResponse, error) {
	sub, err := e.natsConn.SubscribeSync(task.ReplyTo)
	defer sub.Drain()
	if err != nil {
		return ExecuteResponse{}, err
	}

	msg, err := sub.NextMsg(30 * time.Second)
	if err != nil {
		return ExecuteResponse{}, err
	}

	var response ExecuteResponse
	if err := json.Unmarshal(msg.Data, &response); err != nil {
		return ExecuteResponse{}, err
	}

	return response, nil

}

func (e *ExecutorService) lookupFunction(funcID string) (models.Function, error) {
	msg, err := e.natsConn.Request("function.lookup", []byte(funcID), nats.DefaultTimeout)
	if err != nil {
		return models.Function{}, err
	}

	var function models.Function
	if err := json.Unmarshal(msg.Data, &function); err != nil {
		return models.Function{}, err
	}

	return function, nil
}

func (e *ExecutorService) publishExecutionTask(ctx context.Context, task ExecuteTask) error {
	jsonTask, err := json.Marshal(task)
	if err != nil {
		return err
	}

	_, err = e.js.Publish(ctx, "jobs.execute", jsonTask)

	return err
}
