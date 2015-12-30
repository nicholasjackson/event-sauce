package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/alexcesaro/statsd"
	"github.com/facebookgo/inject"
	"github.com/nicholasjackson/event-sauce/data"
	"github.com/nicholasjackson/event-sauce/global"
	"github.com/nicholasjackson/event-sauce/handlers"
	"github.com/nicholasjackson/event-sauce/logging"
	"github.com/nicholasjackson/event-sauce/queue"
	"github.com/nicholasjackson/event-sauce/workers"
)

func main() {
	config := os.Args[1]
	rootfolder := os.Args[2]

	global.LoadConfig(config, rootfolder)

	var wg sync.WaitGroup
	wg.Add(2)

	setupInjection()

	go startApiServer(&wg)
	go startClient(&wg)

	wg.Wait()
}

func setupInjection() {
	err := global.SetupInjection(
		&inject.Object{Value: handlers.HealthHandlerDependencies},
		&inject.Object{Value: handlers.RegisterHandlerDependencies},
		&inject.Object{Value: handlers.EventHandlerDependencies},
		&inject.Object{Value: ClientDeps},
		&inject.Object{Value: createStatsDClient(), Name: "statsd"},
		&inject.Object{Value: createMongoClient(), Name: "dal"},
		&inject.Object{Value: createEventQueueClient(), Name: "eventqueue"},
		&inject.Object{Value: createDeadLetterQueueClient(), Name: "deadletterqueue"},
		&inject.Object{Value: createEventDispatcher(), Name: "eventdispatcher"},
		&inject.Object{Value: createEventQueueWorkerFactory(), Name: "eventqueueworkerfactory"},
	)

	if err != nil {
		panic(fmt.Sprintln("Unable to create injection framework: ", err))
	}

}

func createStatsDClient() logging.StatsD {
	statsDClient, err := statsd.New(global.Config.StatsDServerIP) // reference to a statsd client
	if err != nil {
		panic(fmt.Sprintln("Unable to create StatsD Client: ", err))
	}
	return statsDClient
}

func createMongoClient() *data.MongoDal {
	dal, err := data.New(global.Config.Data.ConnectionString, global.Config.Data.DataBaseName)
	if err != nil {
		panic(fmt.Sprintln("Unable to create DataBase: ", err))
	}
	return dal
}

func createEventQueueClient() *queue.RedisQueue {
	queue, err := queue.New(global.Config.Queue.ConnectionString, global.Config.Queue.EventQueue)
	if err != nil {
		panic(fmt.Sprintln("Unable to create Queue: ", err))
	}
	return queue
}

func createDeadLetterQueueClient() *queue.RedisQueue {
	queue, err := queue.New(global.Config.Queue.ConnectionString, global.Config.Queue.DeadLetterQueue)
	if err != nil {
		panic(fmt.Sprintln("Unable to create Queue: ", err))
	}
	return queue
}

func createEventDispatcher() workers.EventDispatcher {
	return &workers.HTTPEventDispatcher{}
}

func createEventQueueWorkerFactory() workers.WorkerFactory {
	return &workers.EventQueueWorkerFactory{}
}
