package queue

import (
	"encoding/json"
	"fmt"

	"github.com/adjust/rmq"
)

type queuePayload struct {
	MessageName string `json:"message_name"`
	Payload     string `json:"payload"`
}

type RedisQueue struct {
	Queue rmq.Queue
}

func New(connectionString string, queueName string) (*RedisQueue, error) {
	connection := rmq.OpenConnection("my service", "tcp", connectionString, 1)
	taskQueue := connection.OpenQueue(queueName)

	return &RedisQueue{Queue: taskQueue}, nil
}

func (r *RedisQueue) Add(messageName string, payload string) error {
	queuePayload := queuePayload{MessageName: messageName, Payload: payload}
	payloadBytes, err := json.Marshal(queuePayload)
	if err != nil {
		// handle error
		return err
	}
	fmt.Println("Event:", string(payloadBytes))
	r.Queue.PublishBytes(payloadBytes)
	//r.Queue.Publish("bollocks")
	return nil
}
