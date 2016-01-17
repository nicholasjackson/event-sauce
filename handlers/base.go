package handlers

type BaseResponse struct {
	StatusEvent string `json:"status_event"`
}

const (
	GET                 = ".get"
	POST                = ".post"
	PUT                 = ".put"
	DELETE              = ".delete"
	CALLED              = ".called"
	SUCCESS             = ".success"
	PROCESS_REDELIVERY  = ".process_redelivery"
	DELETE_REGISTRATION = ".delete_registration"
	NO_ENDPOINT         = ".no_registered_endpoint"
	HANDLE              = ".handle"
	DISPATCH            = ".dispatch"
	NOT_FOUND           = ".not_found"
	ERROR               = ".server_error"
	INVALID_REQUEST     = ".request.invalid_request"
	BAD_REQUEST         = ".request.bad_request"
	VALID_REQUEST       = ".request.valid"
	INVALID_TOKEN       = ".auth.invalid_token"
	NOT_AUTHORISED      = ".auth.not_authorised"
	TOKEN_OK            = ".auth.token_ok"
	HEALTH_HANDLER      = "eventsauce.health"
	EVENT_HANDLER       = "eventsauce.event"
	REGISTER_HANDLER    = "eventsauce.register"
	EVENT_QUEUE_WORKER  = "eventsauce.event_queue_worker"
	DEAD_LETTER_WORKER  = "eventsauce.dead_letter_worker"
)
