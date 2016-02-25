package workers

import (
	"log"
	"time"

	"github.com/nicholasjackson/sorcery/data"
	"github.com/nicholasjackson/sorcery/entities"
	"github.com/nicholasjackson/sorcery/global"
	"github.com/nicholasjackson/sorcery/handlers"
	"github.com/nicholasjackson/sorcery/logging"
)

type DeadLetterWorker struct {
	eventDispatcher EventDispatcher
	dal             data.Dal
	log             *log.Logger
	statsD          logging.StatsD
}

const DQWTAGNAME = "DeadLetterQueueWorker: "

func NewDeadLetterWorker(
	eventDispatcher EventDispatcher,
	dal data.Dal,
	log *log.Logger,
	statsD logging.StatsD) *DeadLetterWorker {
	return &DeadLetterWorker{eventDispatcher: eventDispatcher, dal: dal, log: log, statsD: statsD}
}

func (w *DeadLetterWorker) HandleItem(item interface{}) error {
	deadLetter := item.(*entities.DeadLetterItem)

	w.log.Printf("%vProcessing event: %v for: %v\n", DQWTAGNAME, deadLetter.Event.EventName, deadLetter.CallbackUrl)
	w.statsD.Increment(handlers.DEAD_LETTER_QUEUE + handlers.WORKER + handlers.HANDLE)

	registration, err := w.dal.GetRegistrationByEventAndCallback(deadLetter.Event.EventName, deadLetter.CallbackUrl)

	if err == nil {
		w.processRegistration(registration, deadLetter)
	} else {
		w.statsD.Increment(handlers.DEAD_LETTER_QUEUE + handlers.WORKER + handlers.NO_ENDPOINT)
	}
	return nil
}

func (w *DeadLetterWorker) processRegistration(registration *entities.Registration, deadLetter *entities.DeadLetterItem) {
	code, _ := w.eventDispatcher.DispatchEvent(&deadLetter.Event, deadLetter.CallbackUrl)

	w.log.Printf(
		"%vEvent dispatched: %v for: %v return code: %v\n",
		DQWTAGNAME,
		deadLetter.Event.EventName,
		deadLetter.CallbackUrl,
		code)

	switch code {
	case 500:
		w.processRedelivery(deadLetter, registration) // the endpoint is unhealthy retry if possible
		break
	case 404:
		w.deleteRegistration(registration) // endpoint does not exist delete endpoint
		break
	case 200:
		w.statsD.Increment(handlers.DEAD_LETTER_QUEUE + handlers.WORKER + handlers.DISPATCH)
		break
	}
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
	w.statsD.Increment(handlers.DEAD_LETTER_QUEUE + handlers.WORKER + handlers.PROCESS_REDELIVERY)

	duration, _ := time.ParseDuration(global.Config.RetryIntervals[deadLetter.FailureCount])
	deadLetter.FailureCount++
	deadLetter.NextRetryDate = deadLetter.NextRetryDate.Add(duration)
	w.dal.UpsertDeadLetterItem(deadLetter)
}

func (w *DeadLetterWorker) deleteRegistration(registration *entities.Registration) {
	w.statsD.Increment(handlers.DEAD_LETTER_QUEUE + handlers.WORKER + handlers.DELETE_REGISTRATION)

	w.dal.DeleteRegistration(registration)
}
