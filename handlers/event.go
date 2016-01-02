package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nicholasjackson/event-sauce/logging"
	"github.com/nicholasjackson/event-sauce/queue"
)

type EventRequest struct {
	EventName string          `json:"event_name"`
	Payload   json.RawMessage `json:"payload"`
}

type EventDependencies struct {
	// statsD interface must use a name type as injection cannot infer ducktypes
	Stats logging.StatsD `inject:"statsd"`
	Queue queue.Queue    `inject:"eventqueue"`
	Log   *log.Logger    `inject:""`
}

var EventHandlerDependencies *EventDependencies = &EventDependencies{}

const EVENT_HANDLER_CALLED = "eventsauce.event_handler.post"
const EHTAGNAME = "EventHandler: "

func EventHandler(rw http.ResponseWriter, r *http.Request) {
	EventHandlerDependencies.Stats.Increment(EVENT_HANDLER_CALLED)
	EventHandlerDependencies.Log.Printf("%vHandler Called POST\n", EHTAGNAME)

	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body)
	request := EventRequest{}

	err := json.Unmarshal(data, &request)
	if err != nil || request.EventName == "" || len(request.Payload) < 1 {
		http.Error(rw, "Invalid request object", http.StatusBadRequest)
		return
	}

	if err = EventHandlerDependencies.Queue.Add(request.EventName, string(request.Payload)); err != nil {
		http.Error(rw, "Error adding item to queue", http.StatusInternalServerError)
		return
	} else {
		var response BaseResponse
		response.StatusEvent = "OK"

		encoder := json.NewEncoder(rw)
		encoder.Encode(&response)

	}
}
