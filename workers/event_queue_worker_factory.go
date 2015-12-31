package workers

import (
	"github.com/nicholasjackson/event-sauce/data"
	"github.com/nicholasjackson/event-sauce/queue"
)

type EventQueueWorkerFactory struct {
	EventDispatcher EventDispatcher `inject:"eventdispatcher"`
	Dal             data.Dal        `inject:"dal"`
	DeadLetterQueue queue.Queue     `inject:"deadletterqueue"`
}

func (f *EventQueueWorkerFactory) Create() Worker {
	return New(f.EventDispatcher, f.Dal, f.DeadLetterQueue)
}
