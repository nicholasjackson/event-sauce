package handlers

type BaseResponse struct {
	StatusEvent string `json:"status_event"`
}

const (
	POST                = ".post"
	PUT                 = ".put"
	DELETE              = ".delete"
	CALLED              = ".called"
	SUCCESS             = ".success"
	NOT_FOUND           = ".not_found"
	ERROR               = ".server_error"
	INVALID_REQUEST     = ".request.invalid_request"
	BAD_REQUEST         = ".request.bad_request"
	VALID_REQUEST       = ".request.valid"
	INVALID_TOKEN       = ".auth.invalid_token"
	NOT_AUTHORISED      = ".auth.not_authorised"
	TOKEN_OK            = ".auth.token_ok"
	EVENT_HANDLER       = "eventsauce.event_handler"
	REGISTER_HANDLER    = "eventsauce.register_handler"
	EVENT_QUEUE_WORKER  = "api_users.event_queue_worker"
	DEAD_LETTER_WORKDER = "api_users.dead_letter_workder"
)
