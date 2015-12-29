package main

import (
	"fmt"
	"net/http"

	"github.com/nicholasjackson/event-sauce/handlers"
)

func startApiServer() {
	http.Handle("/", handlers.GetRouter())

	fmt.Println("Listening for connections on port", 8001)
	http.ListenAndServe(fmt.Sprintf(":%v", 8001), nil)
}
