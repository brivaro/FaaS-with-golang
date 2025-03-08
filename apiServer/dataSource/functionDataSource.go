package dataSource

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"faas/initializers/nclient"
	"faas/models"
	"fmt"
	"log"
	"time"
)

func InsertFunction(function models.Function) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var funcID string
	// Creamos la clave aleatoria que no debe existir en el KVStore
	for {
		funcID = generateRandomKey()
		_, err := nclient.Client.KvFunctions.Get(ctx, funcID)
		// si encuentra func con ese id, repite el for hasta que la key no este en kvFunc
		if err == nil {
			continue
		}
		now := time.Now()

		formattedDate := now.Format("02/01/2006 15:04:05")

		function.ID = funcID
		function.CreatedAt = formattedDate

		funcToJSON, err := json.Marshal(function)
		if err != nil {
			fmt.Print(err)
			return "", errors.New("failed to Marshal function")
		}
		_, err = nclient.Client.KvFunctions.Put(ctx, funcID, funcToJSON)
		if err != nil {
			return "", errors.New("failed to register function")
		}
		break
	}
	return funcID, nil
}

func GetFunctionByID(funcID string) (models.Function, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	entry, err := nclient.Client.KvFunctions.Get(ctx, funcID)
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

func generateRandomKey() string {
	bytes := make([]byte, 16) // Genera 16 bytes aleatorios
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatalf("Error generando clave aleatoria: %v", err)
	}
	return hex.EncodeToString(bytes) // Devuelve la clave en formato hexadecimal
}

func DeleteFunction(functionID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := nclient.Client.KvFunctions.Delete(ctx, functionID)
	if err != nil {
		return fmt.Errorf("failed to delete function: %v", err)
	}

	return nil
}

func GetFunctionsByUsername(userID string) ([]models.Function, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// list all keys of the bucket
	// CURIOSIDAD: si hubiesemos guardado las claves de las funciones como +
	// + "user_function1:{info de la func}" se podrian usar filtros +
	// + se podria usar ...Keys(ctx, jetstream.KeyValueListOptions{ Limit, Prefix, KeysOnly... })
	keys, err := nclient.Client.KvFunctions.Keys(ctx) // sin filtros
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve keys from kv store: %v", err)
	}

	// filter by UserID
	var userFunctions []models.Function
	for _, key := range keys {
		entry, err := nclient.Client.KvFunctions.Get(ctx, key)
		if err != nil {
			log.Printf("Failed to get value for key %s: %v", key, err)
			continue
		}

		var function models.Function
		err = json.Unmarshal(entry.Value(), &function)
		if err != nil {
			log.Printf("Failed to unmarshal function for key %s: %v", key, err)
			continue
		}

		if function.UserID == userID {
			userFunctions = append(userFunctions, function)
		}
	}

	if len(userFunctions) == 0 {
		return nil, errors.New("no functions found for the user")
	}

	return userFunctions, nil
}
