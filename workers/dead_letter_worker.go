package workers

import "github.com/nicholasjackson/event-sauce/data"

type DeadLetterWorker struct {
	eventDispatcher EventDispatcher
	dal             data.Dal
}

func NewDeadLetterWorker(eventDispatcher EventDispatcher, dal data.Dal) *DeadLetterWorker {
	return &DeadLetterWorker{eventDispatcher: eventDispatcher, dal: dal}
}

func (w *DeadLetterWorker) HandleItem(item interface{}) error {
	return nil
}
