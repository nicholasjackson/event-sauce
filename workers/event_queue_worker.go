package workers

import (
	"fmt"

	"github.com/nicholasjackson/event-sauce/data"
	"github.com/nicholasjackson/event-sauce/entities"
	"github.com/nicholasjackson/event-sauce/queue"
)

type EventQueueWorker struct {
	eventDispatcher EventDispatcher
	dal             data.Dal
	deadLetterQueue queue.Queue
}

func New(eventDispatcher EventDispatcher, dal data.Dal, deadLetterQueue queue.Queue) *EventQueueWorker {
	return &EventQueueWorker{eventDispatcher: eventDispatcher, dal: dal, deadLetterQueue: deadLetterQueue}
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
	fmt.Println("Processing Event:", event)

	code, _ := m.eventDispatcher.DispatchEvent(event, registration.CallbackUrl)

	fmt.Println("processEvent: Finshed: ", code)
	switch code {
	case 404:
		fmt.Println("processEvent: Not Found")
		m.deleteRegistration(registration)
		break
	case 200:
		fmt.Println("processEvent: Dispatched OK")
		break
	default:
		fmt.Println("processEvent: Not healthy")
		m.addToDeadLetterQueue(event, registration.CallbackUrl)
		break
	}
}

func (m *EventQueueWorker) deleteRegistration(registration *entities.Registration) {
	m.dal.DeleteRegistration(registration)
}

func (m *EventQueueWorker) addToDeadLetterQueue(event *entities.Event, endpoint string) {
	m.deadLetterQueue.AddEvent(event)
}

func (m *EventQueueWorker) saveEventToStore(event *entities.Event) error {
	eventStore := entities.NewEventStoreItem(*event)
	return m.dal.UpsertEventStore(&eventStore)
}
