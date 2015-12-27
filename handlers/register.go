package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/nicholasjackson/event-sauce/logging"
)

type RegisterRequest struct {
	MessageName string `json:"message_name"`
	HealthUrl   string `json:"health_url"`
	CallbackUrl string `json:"callback_url"`
}

type RegisterDependencies struct {
	// statsD interface must use a name type as injection cannot infer ducktypes
	Stats logging.StatsD `inject:"statsd"`
}

var RegisterHandlerDependencies *RegisterDependencies = &RegisterDependencies{}

const REGISTER_HANDLER_CALLED = "event-sauce.health_handler"

func RegisterHandler(rw http.ResponseWriter, r *http.Request) {
	RegisterHandlerDependencies.Stats.Increment(REGISTER_HANDLER_CALLED)

	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body)
	request := RegisterRequest{}

	err := json.Unmarshal(data, &request)
	if err != nil || request.MessageName == "" || request.HealthUrl == "" || request.CallbackUrl == "" {
		http.Error(rw, "Invalid request object", http.StatusBadRequest)
		return
	}
}
