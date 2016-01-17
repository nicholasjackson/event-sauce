package workers

import (
	"log"

	"github.com/nicholasjackson/event-sauce/data"
	"github.com/transform/api-users/logging"
)

type DeadLetterQueueWorkerFactory struct {
	EventDispatcher EventDispatcher `inject:"eventdispatcher"`
	Dal             data.Dal        `inject:"dal"`
	StatsD          logging.StatsD  `inject:"statsd"`
	Log             *log.Logger     `inject:""`
}

func (f *DeadLetterQueueWorkerFactory) Create() Worker {
	return NewDeadLetterWorker(f.EventDispatcher, f.Dal, f.Log, f.StatsD)
}
