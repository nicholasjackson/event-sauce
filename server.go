package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/nicholasjackson/sorcery/handlers"
)

func startApiServer(wg *sync.WaitGroup) {
	defer wg.Done()

	http.Handle("/", handlers.GetRouter())

	fmt.Println("Listening for connections on port", 8001)
	http.ListenAndServe(fmt.Sprintf(":%v", 8001), nil)
}
