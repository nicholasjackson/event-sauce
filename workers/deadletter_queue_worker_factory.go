package workers

import (
	"log"

	"github.com/nicholasjackson/event-sauce/data"
)

type DeadLetterQueueWorkerFactory struct {
	EventDispatcher EventDispatcher `inject:"eventdispatcher"`
	Dal             data.Dal        `inject:"dal"`
	Log             *log.Logger     `inject:""`
}

func (f *DeadLetterQueueWorkerFactory) Create() Worker {
	return NewDeadLetterWorker(f.EventDispatcher, f.Dal, f.Log)
}
