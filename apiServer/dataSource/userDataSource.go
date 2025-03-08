package dataSource

import (
	"context"
	"encoding/json"
	"errors"
	"faas/initializers/nclient"
	"faas/models"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

func GetUserByUsername(username string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	entry, err := nclient.Client.KvUsers.Get(ctx, username)

	if err != nil {
		return user, err
	}

	err = json.Unmarshal(entry.Value(), &user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func InsertUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userToJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Print(err)
		return errors.New("failed to Marshal user")
	}

	_, err = nclient.Client.KvUsers.Put(ctx, user.Username, userToJSON)
	if err == nats.ErrInvalidKey {
		return errors.New("username is invalid")
	} else if err != nil {
		fmt.Print(err)
		return errors.New("failed to create user")
	}
	return nil
}

func GetAllUsers() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var users []models.User
	keys, err := nclient.Client.KvUsers.Keys(ctx)
	if err != nil {
		return users, err
	}

	for _, key := range keys {
		entry, err := nclient.Client.KvUsers.Get(ctx, key)
		if err != nil {
			return users, err
		}

		var user models.User

		err = json.Unmarshal(entry.Value(), &user)
		if err != nil {
			return users, err
		}

		users = append(users, user)
	}
	return users, nil
}
