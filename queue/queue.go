package queue

type Queue interface {
	Add(message_name string, payload string) error
}
