package queue

import (
	"fmt"
	"time"

	"github.com/nicholasjackson/event-sauce/data"
	"github.com/nicholasjackson/event-sauce/entities"
	"github.com/nicholasjackson/event-sauce/global"
)

type DeadLetterQueue struct {
	Dal data.Dal
}

func NewDeadLetterQueue(dataAccessLayer data.Dal) (*DeadLetterQueue, error) {
	return &DeadLetterQueue{Dal: dataAccessLayer}, nil
}

func (d *DeadLetterQueue) Add(event_name string, payload string) error {
	return nil
}

func (d *DeadLetterQueue) AddEvent(event *entities.Event, callback string) error {
	deadLetter := entities.NewDeadLetterItem(*event)
	duration, _ := time.ParseDuration(global.Config.RetryIntervals[0])

	deadLetter.CallbackUrl = callback
	deadLetter.FailureCount = 1
	deadLetter.FirstFailureDate = time.Now()
	deadLetter.NextRetryDate = deadLetter.FirstFailureDate.Add(duration)

	return d.Dal.UpsertDeadLetterItem(&deadLetter)
}

func (d *DeadLetterQueue) StartConsuming(size int, poll_interval time.Duration, callback func(callbackItem interface{})) {
	for {
		d.runConsumer(size, callback)
		time.Sleep(poll_interval)
	}
}

func (d *DeadLetterQueue) runConsumer(size int, callback func(callbackItem interface{})) {
	deadLetters, err := d.Dal.GetDeadLetterItemsReadyForRetry()

	if deadLetters == nil || err != nil {
		fmt.Println("Nothing to do")
		return
	}

	// remove the retrieved items from the queue to ensure that no other worker picks them up when processing
	// the worker will re-add to the queue in the event of failure
	_ = d.Dal.DeleteDeadLetterItems(deadLetters)

	buffer := make(chan struct{}, size) // make channel same size as consumer size

	for i, letter := range deadLetters {
		buffer <- struct{}{} // add to the channel
		go func(i int, letter *entities.DeadLetterItem) {
			callback(letter)
			<-buffer // release channel
		}(i, letter)
	}
}
