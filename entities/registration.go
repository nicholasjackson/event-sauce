package entities

import (
	"time"

	"labix.org/v2/mgo/bson"
)

type Registration struct {
	Id           bson.ObjectId `bson:"_id"`
	EventName  string        "event_name,omitempty"
	CallbackUrl  string        "callback_url,omitempty"
	CreationDate time.Time     "creation_date,omitempty"
}
