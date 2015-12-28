package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nicholasjackson/event-sauce/data"
	"github.com/nicholasjackson/event-sauce/global"
	"github.com/nicholasjackson/event-sauce/handlers"
	"github.com/nicholasjackson/event-sauce/logging"

	"github.com/alexcesaro/statsd"
	"github.com/facebookgo/inject"
)

func main() {
	config := os.Args[1]
	rootfolder := os.Args[2]

	global.LoadConfig(config, rootfolder)

	setupInjection()
	setupHandlers()
}

func setupHandlers() {
	http.Handle("/", handlers.GetRouter())

	fmt.Println("Listening for connections on port", 8001)
	http.ListenAndServe(fmt.Sprintf(":%v", 8001), nil)
}

func setupInjection() {
	err := global.SetupInjection(
		&inject.Object{Value: handlers.HealthHandlerDependencies},
		&inject.Object{Value: handlers.RegisterHandlerDependencies},
		&inject.Object{Value: createStatsDClient(), Name: "statsd"},
		&inject.Object{Value: createMongoClient(), Name: "dal"},
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
