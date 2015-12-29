package data

import (
	"fmt"
	"log"

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

func (m *MongoDal) GetRegistrationByMessageAndCallback(message string, callback_url string) (*entities.Registration, error) {
	return findRegistration(bson.M{"message_name": message, "callback_url": callback_url}, m)
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

func findRegistration(bson interface{}, d *MongoDal) (*entities.Registration, error) {
	session := d.mainSession.New()
	c := session.DB(d.dataBaseName).C("registrations")
	registration := entities.Registration{}

	err := c.Find(bson).One(&registration)
	if err != nil {
		log.Printf("Find Registration Error: %v\n", err)
		return nil, err
	}
	return &registration, nil
}
