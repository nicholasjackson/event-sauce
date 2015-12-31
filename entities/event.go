package entities

import (
	"time"

	"labix.org/v2/mgo/bson"
)

// this needs separated into data entities and api entities
type Event struct {
	EventName string `json:"event_name" bson:"event_name"`
	Payload   string `json:"payload" bson:"payload"`
}

type EventStoreItem struct {
	Id           bson.ObjectId `bson:"_id"`
	Event        Event         `bson:"event,omitempty"`
	CreationDate time.Time     `bson:"creation_date,omitempty"`
}

type DeadLetterItem struct {
	Id               bson.ObjectId `bson:"_id"`
	Event            Event         `bson:"event,omitempty"`
	CreationDate     time.Time     `bson:"creation_date,omitempty"`
	FirstFailureDate time.Time     `bson:"first_failure_date,omitempty"`
	NextRetryDate    time.Time     `bson:"next_retry_date,omitempty"`
	FailureCount     int           `bson:"failure_count,omitempty"`
	CallbackUrl      string        `bson:"callback_url,omitempty"`
}

func NewEventStoreItem(event Event) EventStoreItem {
	return EventStoreItem{
		Id:           bson.NewObjectId(),
		Event:        event,
		CreationDate: time.Now(),
	}
}

func NewDeadLetterItem(event Event) DeadLetterItem {
	return DeadLetterItem{
		Id:           bson.NewObjectId(),
		Event:        event,
		CreationDate: time.Now(),
	}
}
