package queue

import (
	"time"

	"github.com/nicholasjackson/sorcery/entities"
)

type Queue interface {
	Add(eventName string, payload string) error
	AddEvent(event *entities.Event, callback string) error
	StartConsuming(size int, pollInterval time.Duration, callback func(callbackItem interface{}))
}
