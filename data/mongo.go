package data

import (
	"fmt"
	"log"
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	"github.com/nicholasjackson/event-sauce/entities"
)

type MongoDal struct {
	mainSession  *mgo.Session
	dataBaseName string
}

func New(connectionString string, dataBaseName string) (*MongoDal, error) {
	session, err := mgo.Dial(connectionString)
	if err != nil {
		return nil, err
	}

	return &MongoDal{mainSession: session, dataBaseName: dataBaseName}, nil
}

func (m *MongoDal) GetRegistrationByEventAndCallback(event string, callback_url string) (*entities.Registration, error) {
	query := m.findRegistrations(bson.M{"event_name": event, "callback_url": callback_url})
	registration := entities.Registration{}

	err := query.One(&registration)
	if err != nil {
		log.Printf("Find Registration Error: %v\n", err)
		return nil, err
	}
	return &registration, nil
}

func (m *MongoDal) GetRegistrationsByEvent(event string) ([]*entities.Registration, error) {
	query := m.findRegistrations(bson.M{"event_name": event})
	registrations := []*entities.Registration{}

	err := query.All(&registrations)
	if err != nil {
		log.Printf("Find Registration Error: %v\n", err)
		return nil, err
	}
	return registrations, nil
}

func (m *MongoDal) UpsertRegistration(registration *entities.Registration) error {
	log.Printf("Create new Registration: %v\n", registration)
	registration.Id = bson.NewObjectId()
	session := m.mainSession.New()
	c := session.DB(m.dataBaseName).C("registrations")
	err := c.Insert(registration)

	return err
}

func (m *MongoDal) DeleteRegistration(registration *entities.Registration) error {
	log.Printf("Delete Registration: %v\n", registration)
	session := m.mainSession.New()
	c := session.DB(m.dataBaseName).C("registrations")
	err := c.RemoveId(registration.Id)
	fmt.Println("Error: ", err)

	return err
}

func (m *MongoDal) UpsertEvent(event *entities.Event) error {
	log.Printf("Create new Event: %v\n", event)
	dbEvent := &entities.DBEvent{}
	dbEvent.Id = bson.NewObjectId()
	dbEvent.EventName = event.EventName
	dbEvent.Callback = event.Callback
	dbEvent.Payload = event.Payload
	dbEvent.CreationDate = time.Now()

	session := m.mainSession.New()
	c := session.DB(m.dataBaseName).C("events")
	err := c.Insert(dbEvent)

	return err
}

func (m *MongoDal) findRegistrations(bson interface{}) *mgo.Query {
	session := m.mainSession.New()
	c := session.DB(m.dataBaseName).C("registrations")

	return c.Find(bson)
}
