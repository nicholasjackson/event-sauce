package entities

import (
	"time"

	"labix.org/v2/mgo/bson"
)

type Registration struct {
	Id            bson.ObjectId `bson:"_id"`
	EventName     string        `bson:"event_name,omitempty"`
	CallbackUrl   string        `bson:"callback_url,omitempty"`
	CreationDate  time.Time     `bson:"creation_date,omitempty"`
	FirstFailDate time.Time     `bson:"first_fail_date,omitempty"`
	LastFailDate  time.Time     `bson:"last_fail_date,omitempty"`
}
