package workers

import (
	"log"
	"time"

	"github.com/nicholasjackson/event-sauce/data"
	"github.com/nicholasjackson/event-sauce/entities"
	"github.com/nicholasjackson/event-sauce/global"
)

type DeadLetterWorker struct {
	eventDispatcher EventDispatcher
	dal             data.Dal
	log             *log.Logger
}

const DQWTAGNAME = "DeadLetterQueueWorker: "

func NewDeadLetterWorker(eventDispatcher EventDispatcher, dal data.Dal, log *log.Logger) *DeadLetterWorker {
	return &DeadLetterWorker{eventDispatcher: eventDispatcher, dal: dal, log: log}
}

func (w *DeadLetterWorker) HandleItem(item interface{}) error {
	deadLetter := item.(*entities.DeadLetterItem)

	w.log.Printf("%vProcessing event: %v for: %v\n", DQWTAGNAME, deadLetter.Event.EventName, deadLetter.CallbackUrl)

	registration, err := w.dal.GetRegistrationByEventAndCallback(deadLetter.Event.EventName, deadLetter.CallbackUrl)

	if registration != nil && err == nil {
		code, _ := w.eventDispatcher.DispatchEvent(&deadLetter.Event, deadLetter.CallbackUrl)
		w.log.Printf("%vEvent dispatched: %v for: %v return code: %v\n", DQWTAGNAME, deadLetter.Event.EventName, deadLetter.CallbackUrl, code)
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
		w.log.Printf("%vQueue for redelivery: %v for: %v\n", DQWTAGNAME, deadLetter.Event.EventName, deadLetter.CallbackUrl)
		w.queueForRedelivery(deadLetter)
	} else {
		w.log.Printf("%vDelete registration: %v for: %v\n", DQWTAGNAME, deadLetter.Event.EventName, deadLetter.CallbackUrl)
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
