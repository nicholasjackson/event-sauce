package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nicholasjackson/sorcery/logging"
)

type HealthDependencies struct {
	// statsD interface must use a name type as injection cannot infer ducktypes
	Stats logging.StatsD `inject:"statsd"`
	Log   *log.Logger    `inject:""`
}

var HealthHandlerDependencies *HealthDependencies = &HealthDependencies{}

const HHTAGNAME = "HealthHandler: "

func HealthHandler(rw http.ResponseWriter, r *http.Request) {
	// all HealthHandlerDependencies are automatically created by injection process
	HealthHandlerDependencies.Stats.Increment(HEALTH_HANDLER + GET + CALLED)
	HealthHandlerDependencies.Log.Printf("%vHandler Called GET\n", HHTAGNAME)

	var response BaseResponse
	response.StatusEvent = "OK"

	encoder := json.NewEncoder(rw)
	encoder.Encode(&response)
}
