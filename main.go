package main

import (
	"fmt"
	"os"

	"github.com/alexcesaro/statsd"
	"github.com/facebookgo/inject"
	"github.com/nicholasjackson/event-sauce/data"
	"github.com/nicholasjackson/event-sauce/global"
	"github.com/nicholasjackson/event-sauce/handlers"
	"github.com/nicholasjackson/event-sauce/logging"
	"github.com/nicholasjackson/event-sauce/queue"
)

func main() {
	config := os.Args[1]
	rootfolder := os.Args[2]

	global.LoadConfig(config, rootfolder)

	setupInjection()
	startApiServer()
	startClient()
}

func setupInjection() {
	err := global.SetupInjection(
		&inject.Object{Value: handlers.HealthHandlerDependencies},
		&inject.Object{Value: handlers.RegisterHandlerDependencies},
		&inject.Object{Value: handlers.EventHandlerDependencies},
		&inject.Object{Value: createStatsDClient(), Name: "statsd"},
		&inject.Object{Value: createMongoClient(), Name: "dal"},
		&inject.Object{Value: createRedisClient(), Name: "queue"},
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

func createRedisClient() *queue.RedisQueue {
	queue, err := queue.New(global.Config.Queue.ConnectionString, global.Config.Queue.MessageQueue)
	if err != nil {
		panic(fmt.Sprintln("Unable to create Queue: ", err))
	}
	return queue
}
