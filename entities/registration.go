package entities

import (
	"time"

	"labix.org/v2/mgo/bson"
)

type Registration struct {
	Id           bson.ObjectId `bson:"_id"`
	EventName    string        `bson:"event_name,omitempty"`
	CallbackUrl  string        `bson:"callback_url,omitempty"`
	CreationDate time.Time     `bson:"creation_date,omitempty"`
}

func CreateNewRegistration(event string, callback string) Registration {
	return Registration{
		Id:           bson.NewObjectId(),
		EventName:    event,
		CallbackUrl:  callback,
		CreationDate: time.Now(),
	}
}
