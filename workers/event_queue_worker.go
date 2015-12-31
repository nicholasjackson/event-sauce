package workers

import (
	"fmt"
	"time"

	"github.com/nicholasjackson/event-sauce/data"
	"github.com/nicholasjackson/event-sauce/entities"
	"github.com/nicholasjackson/event-sauce/global"
)

type EventQueueWorker struct {
	eventDispatcher EventDispatcher
	dal             data.Dal
}

func New(eventDispatcher EventDispatcher, dal data.Dal) *EventQueueWorker {
	return &EventQueueWorker{eventDispatcher: eventDispatcher, dal: dal}
}

func (m *EventQueueWorker) HandleEvent(event *entities.Event) error {
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
	deadLetter := entities.NewDeadLetterItem(*event)
	duration, _ := time.ParseDuration(global.Config.RetryIntervals[0])

	deadLetter.FailureCount = 1
	deadLetter.FirstFailureDate = time.Now()
	deadLetter.NextRetryDate = deadLetter.FirstFailureDate.Add(duration)

	m.dal.UpsertDeadLetterItem(&deadLetter)
}

func (m *EventQueueWorker) saveEventToStore(event *entities.Event) error {
	eventStore := entities.NewEventStoreItem(*event)
	return m.dal.UpsertEventStore(&eventStore)
}
