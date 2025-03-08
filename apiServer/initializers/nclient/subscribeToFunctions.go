package nclient

import (
	"context"
	"encoding/json"
	"faas/models"
	"time"

	"github.com/nats-io/nats.go"
)

func SubscribeFunctions() error {
	_, err := Client.Conn.Subscribe("function.lookup", func(msg *nats.Msg) {
		functionID := string(msg.Data)
		function, err := getFunctionByID(functionID)
		if err != nil {
			errorMsg := []byte("can't find function")
			_ = msg.Respond(errorMsg)
			return
		}

		functionJSON, err := json.Marshal(function)
		if err != nil {
			errorMsg := []byte("failed marshaling function")
			_ = msg.Respond(errorMsg)
			return
		}

		_ = msg.Respond(functionJSON)
	})
	if err != nil {
		return err
	}
	return nil
}

func getFunctionByID(funcID string) (models.Function, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	entry, err := Client.KvFunctions.Get(ctx, funcID)
	var function models.Function
	if err != nil {
		return function, err
	}

	err = json.Unmarshal(entry.Value(), &function)
	if err != nil {
		return function, err
	}
	return function, nil
}
