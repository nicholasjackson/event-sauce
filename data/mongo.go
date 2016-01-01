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
	session := m.mainSession.New()
	defer session.Close()

	c := session.DB(m.dataBaseName).C("registrations")
	query := c.Find(bson.M{"event_name": event, "callback_url": callback_url})

	registration := entities.Registration{}

	err := query.One(&registration)
	if err != nil {
		log.Printf("Find Registration Error: %v\n", err)
		return nil, err
	}

	return &registration, nil
}

func (m *MongoDal) GetRegistrationsByEvent(event string) ([]*entities.Registration, error) {
	session := m.mainSession.New()
	defer session.Close()

	c := session.DB(m.dataBaseName).C("registrations")
	query := c.Find(bson.M{"event_name": event})

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

	session := m.mainSession.New()
	defer session.Close()

	c := session.DB(m.dataBaseName).C("registrations")
	_, err := c.UpsertId(registration.Id, registration)

	return err
}

func (m *MongoDal) DeleteRegistration(registration *entities.Registration) error {
	log.Printf("Delete Registration: %v\n", registration)
	session := m.mainSession.New()
	defer session.Close()

	c := session.DB(m.dataBaseName).C("registrations")
	err := c.RemoveId(registration.Id)
	fmt.Println("Error: ", err)

	return err
}

func (m *MongoDal) UpsertEventStore(event *entities.EventStoreItem) error {
	log.Printf("Create new Event: %v\n", event)

	session := m.mainSession.New()
	defer session.Close()

	c := session.DB(m.dataBaseName).C("events")
	_, err := c.UpsertId(event.Id, event)

	return err
}

func (m *MongoDal) UpsertDeadLetterItem(dead *entities.DeadLetterItem) error {
	log.Printf("Create new Dead letter: %v\n", dead)

	session := m.mainSession.New()
	defer session.Close()

	c := session.DB(m.dataBaseName).C("dead_letters")
	_, err := c.UpsertId(dead.Id, dead)

	return err
}

func (m *MongoDal) GetDeadLetterItemsReadyForRetry() ([]*entities.DeadLetterItem, error) {
	deadletters := []*entities.DeadLetterItem{}

	session := m.mainSession.New()
	defer session.Close()

	c := session.DB(m.dataBaseName).C("dead_letters")
	currentTime := time.Now()
	err := c.Find(bson.M{
		"next_retry_date": bson.M{
			"$lte": currentTime,
		},
	}).All(&deadletters)

	if err != nil {
		log.Printf("Find Registration Error: %v\n", err)
		return nil, err
	}

	return deadletters, nil
}

func (m *MongoDal) DeleteDeadLetterItems(dead []*entities.DeadLetterItem) error {
	session := m.mainSession.New()
	defer session.Close()

	c := session.DB(m.dataBaseName).C("dead_letters")

	for _, letter := range dead {
		_ = c.RemoveId(letter.Id)
	}

	return nil
}
