package data

import "github.com/nicholasjackson/event-sauce/entities"

type Dal interface {
	GetRegistrationByMessageAndCallback(message string, callback_url string) (*entities.Registration, error)
	UpsertRegistration(registration *entities.Registration) error
}
