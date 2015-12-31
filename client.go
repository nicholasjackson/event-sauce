package main

import (
	"sync"
	"time"

	"github.com/nicholasjackson/event-sauce/data"
	"github.com/nicholasjackson/event-sauce/entities"
	"github.com/nicholasjackson/event-sauce/logging"
	"github.com/nicholasjackson/event-sauce/queue"
	"github.com/nicholasjackson/event-sauce/workers"
)

type ClientDependencies struct {
	// statsD interface must use a name type as injection cannot infer ducktypes
	Stats         logging.StatsD        `inject:"statsd"`
	DataStore     data.Dal              `inject:"dal"`
	EventQueue    queue.Queue           `inject:"eventqueue"`
	WorkerFactory workers.WorkerFactory `inject:"eventqueueworkerfactory"`
}

var ClientDeps *ClientDependencies = &ClientDependencies{}

const CLIENT_STARTED = "eventsauce.client.started"

func startClient(wg *sync.WaitGroup) {
	defer wg.Done()

	ClientDeps.Stats.Increment(CLIENT_STARTED)

	ClientDeps.EventQueue.StartConsuming(10, time.Second, func(event *entities.Event) {
		worker := ClientDeps.WorkerFactory.Create()
		worker.HandleEvent(event)
	})
}
