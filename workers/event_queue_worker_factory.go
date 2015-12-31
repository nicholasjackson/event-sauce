package workers

import "github.com/nicholasjackson/event-sauce/data"

type EventQueueWorkerFactory struct {
	EventDispatcher EventDispatcher `inject:"eventdispatcher"`
	Dal             data.Dal        `inject:"dal"`
}

func (f *EventQueueWorkerFactory) Create() Worker {
	return New(f.EventDispatcher, f.Dal)
}
