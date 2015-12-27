package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nicholasjackson/event-sauce/global"
	"github.com/nicholasjackson/event-sauce/handlers"

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

	statsDClient, err := statsd.New(global.Config.StatsDServerIP) // reference to a statsd client
	if err != nil {
		panic(fmt.Sprintln("Unable to create StatsD Client: ", err))
	}

	err = global.SetupInjection(
		&inject.Object{Value: handlers.HealthHandlerDependencies},
		&inject.Object{Value: statsDClient, Name: "statsd"},
	)
	if err != nil {
		panic(fmt.Sprintln("Unable to create injection framework: ", err))
	}

}
