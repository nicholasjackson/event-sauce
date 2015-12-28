package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/nicholasjackson/event-sauce/data"
	"github.com/nicholasjackson/event-sauce/entities"
	"github.com/nicholasjackson/event-sauce/logging"
)

type RegisterRequest struct {
	MessageName string `json:"message_name"`
	CallbackUrl string `json:"callback_url"`
}

type RegisterDependencies struct {
	// statsD interface must use a name type as injection cannot infer ducktypes
	Stats logging.StatsD `inject:"statsd"`
	Dal   data.Dal       `inject:"dal"`
}

var RegisterHandlerDependencies *RegisterDependencies = &RegisterDependencies{}

const REGISTER_HANDLER_CALLED = "event-sauce.health_handler"

func RegisterHandler(rw http.ResponseWriter, r *http.Request) {
	RegisterHandlerDependencies.Stats.Increment(REGISTER_HANDLER_CALLED)

	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body)
	request := RegisterRequest{}

	err := json.Unmarshal(data, &request)
	if err != nil || request.MessageName == "" || request.CallbackUrl == "" {
		http.Error(rw, "Invalid request object", http.StatusBadRequest)
		return
	}

	if r, _ := RegisterHandlerDependencies.Dal.GetRegistrationByMessageAndCallback(
		request.MessageName, request.CallbackUrl); r == nil {
		registration := &entities.Registration{}
		registration.MessageName = request.MessageName
		registration.CallbackUrl = request.CallbackUrl
		_ = RegisterHandlerDependencies.Dal.UpsertRegistration(registration)
	} else {
		http.Error(rw, "Registration not modified", http.StatusNotModified)
	}

	var response BaseResponse
	response.StatusMessage = "OK"

	encoder := json.NewEncoder(rw)
	encoder.Encode(&response)
}
