package workers

import (
	"time"

	"github.com/nicholasjackson/event-sauce/data"
	"github.com/nicholasjackson/event-sauce/entities"
	"github.com/nicholasjackson/event-sauce/global"
)

type DeadLetterWorker struct {
	eventDispatcher EventDispatcher
	dal             data.Dal
}

func NewDeadLetterWorker(eventDispatcher EventDispatcher, dal data.Dal) *DeadLetterWorker {
	return &DeadLetterWorker{eventDispatcher: eventDispatcher, dal: dal}
}

func (w *DeadLetterWorker) HandleItem(item interface{}) error {
	deadLetter := item.(*entities.DeadLetterItem)
	registration, err := w.dal.GetRegistrationByEventAndCallback(deadLetter.Event.EventName, deadLetter.CallbackUrl)

	if registration != nil && err == nil {
		code, _ := w.eventDispatcher.DispatchEvent(&deadLetter.Event, deadLetter.CallbackUrl)
		switch code {
		case 500:
			w.processRedelivery(deadLetter, registration) // the endpoint is unhealthy retry if possible
			break
		case 404:
			w.deleteRegistration(registration) // endpoint does not exist delete endpoint
			break
		}
	}
	return nil
}

func (w *DeadLetterWorker) processRedelivery(deadLetter *entities.DeadLetterItem, registration *entities.Registration) {
	if deadLetter.FailureCount < len(global.Config.RetryIntervals) {
		w.queueForRedelivery(deadLetter)
	} else {
		w.deleteRegistration(registration)
	}
}

func (w *DeadLetterWorker) queueForRedelivery(deadLetter *entities.DeadLetterItem) {
	duration, _ := time.ParseDuration(global.Config.RetryIntervals[deadLetter.FailureCount])
	deadLetter.FailureCount++
	deadLetter.NextRetryDate = deadLetter.NextRetryDate.Add(duration)
	w.dal.UpsertDeadLetterItem(deadLetter)
}

func (w *DeadLetterWorker) deleteRegistration(registration *entities.Registration) {
	w.dal.DeleteRegistration(registration)
}
