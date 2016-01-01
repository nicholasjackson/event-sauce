package queue

import (
	"time"

	"github.com/nicholasjackson/event-sauce/entities"
)

type Queue interface {
	Add(event_name string, payload string) error
	AddEvent(event *entities.Event, callback string) error
	StartConsuming(size int, poll_interval time.Duration, callback func(callbackItem interface{}))
}
