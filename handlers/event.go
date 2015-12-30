package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/nicholasjackson/event-sauce/logging"
	"github.com/nicholasjackson/event-sauce/queue"
)

type EventRequest struct {
	MessageName string          `json:"message_name"`
	Payload     json.RawMessage `json:"payload"`
}

type EventDependencies struct {
	// statsD interface must use a name type as injection cannot infer ducktypes
	Stats logging.StatsD `inject:"statsd"`
	Queue queue.Queue    `inject:"eventqueue"`
}

var EventHandlerDependencies *EventDependencies = &EventDependencies{}

const EVENT_HANDLER_CALLED = "eventsauce.event_handler.new"

func EventHandler(rw http.ResponseWriter, r *http.Request) {
	EventHandlerDependencies.Stats.Increment(EVENT_HANDLER_CALLED)

	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body)
	request := EventRequest{}

	err := json.Unmarshal(data, &request)
	if err != nil || request.MessageName == "" || len(request.Payload) < 1 {
		http.Error(rw, "Invalid request object", http.StatusBadRequest)
		return
	}

	if err = EventHandlerDependencies.Queue.Add(request.MessageName, string(request.Payload)); err != nil {
		http.Error(rw, "Error adding item to queue", http.StatusInternalServerError)
		return
	} else {
		var response BaseResponse
		response.StatusMessage = "OK"

		encoder := json.NewEncoder(rw)
		encoder.Encode(&response)

	}
}
