package entities

import (
	"time"

	"labix.org/v2/mgo/bson"
)

// this needs separated into data entities and api entities
type Event struct {
	Id          string `json:"id"`
	MessageName string `json:"message_name"`
	Payload     string `json:"payload"`
	Callback    string `json:"callback"`
}

type DBEvent struct {
	Id           bson.ObjectId `bson:"_id"`
	MessageName  string        "message_name,omitempty"
	Payload      string        "payload,omitempty"
	Callback     string        "callback,omitempty"
	CreationDate time.Time     "creation_date,omitempty"
}
