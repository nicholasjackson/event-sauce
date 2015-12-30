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

func New(eventDispatcher EventDispatcher, dal data.Dal, queue queue.Queue) *EventQueueWorker {
	return &EventQueueWorker{eventDispatcher: eventDispatcher, dal: dal, deadLetterQueue: queue}
}

func (m *EventQueueWorker) HandleMessage(event *entities.Event) error {
	_ = m.saveEventToStore(event)

	registrations, _ := m.dal.GetRegistrationsByMessage(event.MessageName)

	for _, registration := range registrations {
		m.processEvent(event, registration)
	}

	return nil
}

func (m *EventQueueWorker) processEvent(event *entities.Event, registration *entities.Registration) {
	fmt.Println("Processing Event:", event)

	code, _ := m.eventDispatcher.DispatchMessage(event, registration.CallbackUrl)

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
	event.Callback = endpoint
	m.deadLetterQueue.AddEvent(event)
}

func (m *EventQueueWorker) saveEventToStore(event *entities.Event) error {
	return m.dal.UpsertEvent(event)
}
