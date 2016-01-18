package queue

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/adjust/rmq"
	"github.com/nicholasjackson/event-sauce/entities"
)

type RedisQueue struct {
	Queue    rmq.Queue
	name     string
	callback func(callbackItem interface{})
}

func NewRedisQueue(connectionString string, queueName string) (*RedisQueue, error) {
	connection := rmq.OpenConnection("my service", "tcp", connectionString, 1)
	taskQueue := connection.OpenQueue(queueName)

	return &RedisQueue{Queue: taskQueue, name: queueName}, nil
}

func (r *RedisQueue) Add(eventName string, payload string) error {
	queuePayload := entities.Event{EventName: eventName, Payload: payload}

	return r.AddEvent(&queuePayload, "")
}

func (r *RedisQueue) AddEvent(event *entities.Event, callback string) error {
	payloadBytes, err := json.Marshal(event)
	if err != nil {
		// handle error
		return err
	}

	fmt.Println("Add event to queue:", string(payloadBytes))
	r.Queue.PublishBytes(payloadBytes)

	return nil
}

func (r *RedisQueue) StartConsuming(size int, pollInterval time.Duration, callback func(callbackItem interface{})) {
	r.callback = callback
	r.Queue.StartConsuming(size, pollInterval)
	r.Queue.AddConsumer("RedisQueue_"+r.name, r)
}

// Interface from rmq
func (r *RedisQueue) Consume(delivery rmq.Delivery) {
	fmt.Println("Got event from queue:", delivery.Payload())

	event := entities.Event{}

	if err := json.Unmarshal([]byte(delivery.Payload()), &event); err != nil {
		fmt.Println("Error consuming event, unable to deserialise event")
		// handle error
		delivery.Reject()
		return
	}
	delivery.Ack()

	r.callback(&event)
}
