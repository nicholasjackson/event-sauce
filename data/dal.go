package data

import "github.com/nicholasjackson/sorcery/entities"

type Dal interface {
	GetRegistrationsByEvent(event string) ([]*entities.Registration, error)
	GetRegistrationByEventAndCallback(event string, callback_url string) (*entities.Registration, error)
	UpsertRegistration(registration *entities.Registration) error
	DeleteRegistration(registration *entities.Registration) error

	UpsertEventStore(event *entities.EventStoreItem) error

	UpsertDeadLetterItem(dead *entities.DeadLetterItem) error
	GetDeadLetterItemsReadyForRetry() ([]*entities.DeadLetterItem, error)
	DeleteDeadLetterItems(dead []*entities.DeadLetterItem) error
}
