package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nicholasjackson/sorcery/data"
	"github.com/nicholasjackson/sorcery/entities"
	"github.com/nicholasjackson/sorcery/logging"
)

type RegisterRequest struct {
	EventName   string `json:"event_name"`
	CallbackUrl string `json:"callback_url"`
}

type RegisterDependencies struct {
	// statsD interface must use a name type as injection cannot infer ducktypes
	Stats logging.StatsD `inject:"statsd"`
	Dal   data.Dal       `inject:"dal"`
	Log   *log.Logger    `inject:""`
}

var RegisterHandlerDependencies *RegisterDependencies = &RegisterDependencies{}

const RHTAGNAME = "RegisterHandler: "

func RegisterCreateHandler(rw http.ResponseWriter, r *http.Request) {
	RegisterHandlerDependencies.Stats.Increment(REGISTER_HANDLER + POST + CALLED)
	RegisterHandlerDependencies.Log.Printf("%vHandler Called POST\n", RHTAGNAME)

	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body)
	request := RegisterRequest{}

	err := json.Unmarshal(data, &request)
	if err != nil || request.EventName == "" || request.CallbackUrl == "" {
		RegisterHandlerDependencies.Stats.Increment(REGISTER_HANDLER + POST + BAD_REQUEST)
		http.Error(rw, "Invalid request object", http.StatusBadRequest)
		return
	}

	if r, _ := RegisterHandlerDependencies.Dal.GetRegistrationByEventAndCallback(
		request.EventName, request.CallbackUrl); r == nil {
		registration := entities.CreateNewRegistration(request.EventName, request.CallbackUrl)
		_ = RegisterHandlerDependencies.Dal.UpsertRegistration(&registration)
	} else {
		RegisterHandlerDependencies.Stats.Increment(REGISTER_HANDLER + POST + NOT_FOUND)
		http.Error(rw, "Registration not modified", http.StatusNotModified)
		return
	}

	RegisterHandlerDependencies.Stats.Increment(REGISTER_HANDLER + POST + SUCCESS)
	var response BaseResponse
	response.StatusEvent = "OK"

	encoder := json.NewEncoder(rw)
	encoder.Encode(&response)
}

func RegisterDeleteHandler(rw http.ResponseWriter, r *http.Request) {
	RegisterHandlerDependencies.Stats.Increment(REGISTER_HANDLER + DELETE + CALLED)
	RegisterHandlerDependencies.Log.Printf("%vHandler Called DELETE\n", RHTAGNAME)

	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body)
	request := RegisterRequest{}

	err := json.Unmarshal(data, &request)
	if err != nil || request.EventName == "" || request.CallbackUrl == "" {
		RegisterHandlerDependencies.Stats.Increment(REGISTER_HANDLER + DELETE + BAD_REQUEST)
		http.Error(rw, "Invalid request object", http.StatusBadRequest)
		return
	}

	if r, _ := RegisterHandlerDependencies.Dal.GetRegistrationByEventAndCallback(
		request.EventName, request.CallbackUrl); r != nil {
		if err = RegisterHandlerDependencies.Dal.DeleteRegistration(r); err != nil {
			RegisterHandlerDependencies.Stats.Increment(REGISTER_HANDLER + DELETE + ERROR)
			http.Error(rw, "Unable to delete request object", http.StatusInternalServerError)
			return
		}

		RegisterHandlerDependencies.Stats.Increment(REGISTER_HANDLER + DELETE + SUCCESS)
		var response BaseResponse
		response.StatusEvent = "OK"

		encoder := json.NewEncoder(rw)
		encoder.Encode(&response)
	} else {
		RegisterHandlerDependencies.Stats.Increment(REGISTER_HANDLER + DELETE + NOT_FOUND)
		http.Error(rw, "Registration not found", http.StatusNotModified)
		return
	}
}
