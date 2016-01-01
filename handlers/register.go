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
	EventName   string `json:"event_name"`
	CallbackUrl string `json:"callback_url"`
}

type RegisterDependencies struct {
	// statsD interface must use a name type as injection cannot infer ducktypes
	Stats logging.StatsD `inject:"statsd"`
	Dal   data.Dal       `inject:"dal"`
}

var RegisterHandlerDependencies *RegisterDependencies = &RegisterDependencies{}

const REGISTER_CREATE_HANDLER_CALLED = "eventsauce.register_handler.create"
const REGISTER_DELETE_HANDLER_CALLED = "eventsauce.register_handler.delete"

func RegisterCreateHandler(rw http.ResponseWriter, r *http.Request) {
	RegisterHandlerDependencies.Stats.Increment(REGISTER_CREATE_HANDLER_CALLED)

	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body)
	request := RegisterRequest{}

	err := json.Unmarshal(data, &request)
	if err != nil || request.EventName == "" || request.CallbackUrl == "" {
		http.Error(rw, "Invalid request object", http.StatusBadRequest)
		return
	}

	if r, _ := RegisterHandlerDependencies.Dal.GetRegistrationByEventAndCallback(
		request.EventName, request.CallbackUrl); r == nil {
		registration := entities.CreateNewRegistration(request.EventName, request.CallbackUrl)
		_ = RegisterHandlerDependencies.Dal.UpsertRegistration(&registration)
	} else {
		http.Error(rw, "Registration not modified", http.StatusNotModified)
	}

	var response BaseResponse
	response.StatusEvent = "OK"

	encoder := json.NewEncoder(rw)
	encoder.Encode(&response)
}

func RegisterDeleteHandler(rw http.ResponseWriter, r *http.Request) {
	RegisterHandlerDependencies.Stats.Increment(REGISTER_DELETE_HANDLER_CALLED)

	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body)
	request := RegisterRequest{}

	err := json.Unmarshal(data, &request)
	if err != nil || request.EventName == "" || request.CallbackUrl == "" {
		http.Error(rw, "Invalid request object", http.StatusBadRequest)
		return
	}

	if r, _ := RegisterHandlerDependencies.Dal.GetRegistrationByEventAndCallback(
		request.EventName, request.CallbackUrl); r != nil {
		if err = RegisterHandlerDependencies.Dal.DeleteRegistration(r); err != nil {
			http.Error(rw, "Unable to delete request object", http.StatusInternalServerError)
			return
		}

		var response BaseResponse
		response.StatusEvent = "OK"

		encoder := json.NewEncoder(rw)
		encoder.Encode(&response)
	} else {
		http.Error(rw, "Registration not found", http.StatusNotFound)
	}
}
