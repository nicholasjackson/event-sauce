package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nicholasjackson/event-sauce/logging"
)

type HealthDependencies struct {
	// statsD interface must use a name type as injection cannot infer ducktypes
	Stats logging.StatsD `inject:"statsd"`
}

var HealthHandlerDependencies *HealthDependencies = &HealthDependencies{}

const HEALTH_HANDLER_CALLED = "event-sauce.health_handler"

func HealthHandler(rw http.ResponseWriter, r *http.Request) {
	// all HealthHandlerDependencies are automatically created by injection process
	HealthHandlerDependencies.Stats.Increment(HEALTH_HANDLER_CALLED)

	var response BaseResponse
	response.StatusEvent = "OK"

	encoder := json.NewEncoder(rw)
	encoder.Encode(&response)
}
