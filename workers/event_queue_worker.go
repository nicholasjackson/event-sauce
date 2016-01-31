package workers

import (
	"log"

	"github.com/nicholasjackson/sorcery/data"
	"github.com/nicholasjackson/sorcery/entities"
	"github.com/nicholasjackson/sorcery/handlers"
	"github.com/nicholasjackson/sorcery/queue"
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

func New(
	eventDispatcher EventDispatcher,
	dal data.Dal,
	deadLetterQueue queue.Queue,
	log *log.Logger,
	statsD logging.StatsD) *EventQueueWorker {
	return &EventQueueWorker{
		eventDispatcher: eventDispatcher,
		dal:             dal,
		deadLetterQueue: deadLetterQueue,
		log:             log,
		statsD:          statsD,
	}
}

func (m *EventQueueWorker) HandleItem(item interface{}) error {
	m.statsD.Increment(handlers.EVENT_QUEUE + handlers.WORKER + handlers.HANDLE)

	event := item.(*entities.Event)
	_ = m.saveEventToStore(event)

	registrations, _ := m.dal.GetRegistrationsByEvent(event.EventName)
	m.processRegistrations(registrations, event)

	return nil
}

func (m *EventQueueWorker) processRegistrations(registrations []*entities.Registration, event *entities.Event) {
	if len(registrations) < 1 {
		m.log.Printf("%vNo registered endpoint for: %v\n", EQWTAGNAME, event.EventName)
		m.statsD.Increment(handlers.EVENT_QUEUE + handlers.WORKER + handlers.NO_ENDPOINT)
	} else {
		for _, registration := range registrations {
			m.processEvent(event, registration)
		}
	}
}

func (m *EventQueueWorker) processEvent(event *entities.Event, registration *entities.Registration) {
	m.log.Printf("%vProcessing event: %v for: %v\n", EQWTAGNAME, event.EventName, registration.CallbackUrl)

	code, _ := m.eventDispatcher.DispatchEvent(event, registration.CallbackUrl)

	m.log.Printf(
		"%vEvent dispatched: %v for: %v return code: %v\n",
		EQWTAGNAME,
		event.EventName,
		registration.CallbackUrl,
		code)

	switch code {
	case 404:
		m.deleteRegistration(registration)
		break
	case 200:
		m.log.Println(EQWTAGNAME, "processEvent: Dispatched OK")
		m.statsD.Increment(handlers.EVENT_QUEUE + handlers.WORKER + handlers.DISPATCH)
		break
	default:
		m.addToDeadLetterQueue(event, registration.CallbackUrl)
		break
	}
}

func (m *EventQueueWorker) deleteRegistration(registration *entities.Registration) {
	m.log.Printf("%vDelete registration: %v for: %v\n", EQWTAGNAME, registration.EventName, registration.CallbackUrl)
	m.statsD.Increment(handlers.EVENT_QUEUE + handlers.WORKER + handlers.DELETE_REGISTRATION)

	m.dal.DeleteRegistration(registration)
}

func (m *EventQueueWorker) addToDeadLetterQueue(event *entities.Event, endpoint string) {
	m.log.Printf("%vQueue for redelivery: %v for: %v\n", EQWTAGNAME, event.EventName, endpoint)
	m.statsD.Increment(handlers.EVENT_QUEUE + handlers.WORKER + handlers.PROCESS_REDELIVERY)

	m.deadLetterQueue.AddEvent(event, endpoint)
}

func (m *EventQueueWorker) saveEventToStore(event *entities.Event) error {
	m.log.Printf("%vSave to event store: %v\n", EQWTAGNAME, event.EventName)

	eventStore := entities.NewEventStoreItem(*event)
	return m.dal.UpsertEventStore(&eventStore)
}
