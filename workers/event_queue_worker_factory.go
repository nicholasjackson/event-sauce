package workers

import (
	"log"

	"github.com/nicholasjackson/event-sauce/data"
	"github.com/nicholasjackson/event-sauce/queue"
)

type EventQueueWorkerFactory struct {
	EventDispatcher EventDispatcher `inject:"eventdispatcher"`
	Dal             data.Dal        `inject:"dal"`
	DeadLetterQueue queue.Queue     `inject:"deadletterqueue"`
	Log             *log.Logger     `inject:""`
}

func (f *EventQueueWorkerFactory) Create() Worker {
	return New(f.EventDispatcher, f.Dal, f.DeadLetterQueue, f.Log)
}
