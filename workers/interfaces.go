package workers

import "github.com/nicholasjackson/event-sauce/entities"

type Worker interface {
	HandleMessage(event *entities.Event) error
}

type WorkerFactory interface {
	Create() Worker
}

type EventDispatcher interface {
	DispatchMessage(event *entities.Event, endpoint string) (int, error)
}
