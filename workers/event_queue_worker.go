package workers

import (
	"log"

	"github.com/nicholasjackson/event-sauce/data"
	"github.com/nicholasjackson/event-sauce/entities"
	"github.com/nicholasjackson/event-sauce/queue"
)

type EventQueueWorker struct {
	eventDispatcher EventDispatcher
	dal             data.Dal
	deadLetterQueue queue.Queue
	log             *log.Logger
}

const EQWTAGNAME = "EventQueueWorker: "

func New(eventDispatcher EventDispatcher, dal data.Dal, deadLetterQueue queue.Queue, log *log.Logger) *EventQueueWorker {
	return &EventQueueWorker{eventDispatcher: eventDispatcher, dal: dal, deadLetterQueue: deadLetterQueue, log: log}
}

func (m *EventQueueWorker) HandleItem(item interface{}) error {
	event := item.(*entities.Event)
	_ = m.saveEventToStore(event)

	registrations, _ := m.dal.GetRegistrationsByEvent(event.EventName)

	for _, registration := range registrations {
		m.processEvent(event, registration)
	}

	return nil
}

func (m *EventQueueWorker) processEvent(event *entities.Event, registration *entities.Registration) {
	m.log.Printf("%vProcessing event: %v for: %v\n", EQWTAGNAME, event.EventName, registration.CallbackUrl)

	code, _ := m.eventDispatcher.DispatchEvent(event, registration.CallbackUrl)
	m.log.Printf("%vEvent dispatched: %v for: %v return code: %v\n", EQWTAGNAME, event.EventName, registration.CallbackUrl, code)

	switch code {
	case 404:
		m.log.Println(EQWTAGNAME, "processEvent: Not Found")
		m.deleteRegistration(registration)
		break
	case 200:
		m.log.Println(EQWTAGNAME, "processEvent: Dispatched OK")
		break
	default:
		m.log.Println(EQWTAGNAME, "processEvent: Not healthy")
		m.addToDeadLetterQueue(event, registration.CallbackUrl)
		break
	}
}

func (m *EventQueueWorker) deleteRegistration(registration *entities.Registration) {
	m.log.Printf("%vDelete registration: %v for: %v\n", EQWTAGNAME, registration.EventName, registration.CallbackUrl)
	m.dal.DeleteRegistration(registration)
}

func (m *EventQueueWorker) addToDeadLetterQueue(event *entities.Event, endpoint string) {
	m.log.Printf("%vQueue for redelivery: %v for: %v\n", EQWTAGNAME, event.EventName, endpoint)
	m.deadLetterQueue.AddEvent(event, endpoint)
}

func (m *EventQueueWorker) saveEventToStore(event *entities.Event) error {
	m.log.Printf("%vSave to event store: %v\n", EQWTAGNAME, event.EventName)
	eventStore := entities.NewEventStoreItem(*event)
	return m.dal.UpsertEventStore(&eventStore)
}
