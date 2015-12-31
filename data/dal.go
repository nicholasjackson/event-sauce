package data

import "github.com/nicholasjackson/event-sauce/entities"

type Dal interface {
	GetRegistrationsByEvent(event string) ([]*entities.Registration, error)
	GetRegistrationByEventAndCallback(event string, callback_url string) (*entities.Registration, error)
	UpsertRegistration(registration *entities.Registration) error
	DeleteRegistration(registration *entities.Registration) error

	UpsertEvent(event *entities.Event) error
}
