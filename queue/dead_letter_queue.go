package queue

import (
	"time"

	"github.com/nicholasjackson/event-sauce/data"
	"github.com/nicholasjackson/event-sauce/entities"
	"github.com/nicholasjackson/event-sauce/global"
)

type DeadLetterQueue struct {
	Dal data.Dal `inject:"dal"`
}

func (d *DeadLetterQueue) Add(event_name string, payload string) error {
	return nil
}

func (d *DeadLetterQueue) AddEvent(event *entities.Event) error {
	deadLetter := entities.NewDeadLetterItem(*event)
	duration, _ := time.ParseDuration(global.Config.RetryIntervals[0])

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
	deadLetters, _ := d.Dal.GetDeadLetterItemsReadyForRetry()
	buffer := make(chan struct{}, size)

	for i, letter := range deadLetters {
		buffer <- struct{}{} // add to the channel
		go func(i int, letter *entities.DeadLetterItem) {
			callback(letter)
			<-buffer // release channel
		}(i, letter)
	}
}
