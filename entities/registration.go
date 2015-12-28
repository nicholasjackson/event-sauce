package entities

import (
	"time"

	"labix.org/v2/mgo/bson"
)

type Registration struct {
	Id           bson.ObjectId `bson:"_id"`
	MessageName  string        "message_name,omitempty"
	CallbackUrl  string        "callback_url,omitempty"
	CreationDate time.Time     "creation_date,omitempty"
}
