package workers

import "github.com/nicholasjackson/event-sauce/entities"

type Worker interface {
	HandleItem(item interface{}) error
}

type WorkerFactory interface {
	Create() Worker
}

type EventDispatcher interface {
	DispatchEvent(event *entities.Event, endpoint string) (int, error)
}
