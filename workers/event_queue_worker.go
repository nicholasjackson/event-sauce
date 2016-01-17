package workers

import (
	"log"

	"github.com/nicholasjackson/event-sauce/data"
	"github.com/nicholasjackson/event-sauce/entities"
	"github.com/nicholasjackson/event-sauce/handlers"
	"github.com/nicholasjackson/event-sauce/queue"
	"github.com/transform/api-users/logging"
)

type EventQueueWorker struct {
	eventDispatcher EventDispatcher
	dal             data.Dal
	deadLetterQueue queue.Queue
	log             *log.Logger
	statsD          logging.StatsD
}

const EQWTAGNAME = "EventQueueWorker: "

func New(eventDispatcher EventDispatcher, dal data.Dal, deadLetterQueue queue.Queue, log *log.Logger, statsD logging.StatsD) *EventQueueWorker {
	return &EventQueueWorker{eventDispatcher: eventDispatcher, dal: dal, deadLetterQueue: deadLetterQueue, log: log, statsD: statsD}
}

func (m *EventQueueWorker) HandleItem(item interface{}) error {
	m.statsD.Increment(handlers.EVENT_QUEUE + handlers.WORKER + handlers.HANDLE)

	event := item.(*entities.Event)
	_ = m.saveEventToStore(event)

	registrations, _ := m.dal.GetRegistrationsByEvent(event.EventName)

	if len(registrations) < 1 {
		m.log.Printf("%vNo registered endpoint for: %v\n", EQWTAGNAME, event.EventName)
		m.statsD.Increment(handlers.EVENT_QUEUE + handlers.WORKER + handlers.NO_ENDPOINT)
		return nil
	}

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
		m.statsD.Increment(handlers.EVENT_QUEUE + handlers.WORKER + handlers.DELETE_REGISTRATION)
		break
	case 200:
		m.log.Println(EQWTAGNAME, "processEvent: Dispatched OK")
		m.statsD.Increment(handlers.EVENT_QUEUE + handlers.WORKER + handlers.DISPATCH)
		break
	default:
		m.log.Println(EQWTAGNAME, "processEvent: Not healthy")
		m.addToDeadLetterQueue(event, registration.CallbackUrl)
		m.statsD.Increment(handlers.EVENT_QUEUE + handlers.WORKER + handlers.PROCESS_REDELIVERY)
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
