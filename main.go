package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/alexcesaro/statsd"
	"github.com/facebookgo/inject"
	"github.com/nicholasjackson/sorcery/data"
	"github.com/nicholasjackson/sorcery/global"
	"github.com/nicholasjackson/sorcery/handlers"
	"github.com/nicholasjackson/sorcery/logging"
	"github.com/nicholasjackson/sorcery/queue"
	"github.com/nicholasjackson/sorcery/workers"
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
		&inject.Object{Value: createLogger()},
		&inject.Object{Value: createStatsDClient(), Name: "statsd"},
		&inject.Object{Value: createMongoClient(), Name: "dal"},
		&inject.Object{Value: createEventQueueClient(), Name: "eventqueue"},
		&inject.Object{Value: createDeadLetterQueueClient(), Name: "deadletterqueue"},
		&inject.Object{Value: createEventDispatcher(), Name: "eventdispatcher"},
		&inject.Object{Value: createEventQueueWorkerFactory(), Name: "eventqueueworkerfactory"},
		&inject.Object{Value: createDeadLetterQueueWorkerFactory(), Name: "deadletterqueueworkerfactory"},
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

func createEventDispatcher() workers.EventDispatcher {
	return &workers.HTTPEventDispatcher{}
}

func createEventQueueClient() *queue.RedisQueue {
	queue, err := queue.NewRedisQueue(global.Config.Queue.ConnectionString, global.Config.Queue.EventQueue)
	if err != nil {
		panic(fmt.Sprintln("Unable to create Queue: ", err))
	}
	return queue
}

func createEventQueueWorkerFactory() workers.WorkerFactory {
	return &workers.EventQueueWorkerFactory{}
}

func createDeadLetterQueueClient() *queue.DeadLetterQueue {
	queue, err := queue.NewDeadLetterQueue(createMongoClient())
	if err != nil {
		panic(fmt.Sprintln("Unable to create Queue: ", err))
	}
	return queue
}

func createDeadLetterQueueWorkerFactory() workers.WorkerFactory {
	return &workers.DeadLetterQueueWorkerFactory{}
}

func createLogger() *log.Logger {
	return log.New(os.Stdout, "EventSauce: ", log.Lshortfile)
}
