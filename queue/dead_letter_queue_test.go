package queue

import (
	"testing"
	"time"

	"github.com/nicholasjackson/event-sauce/entities"
	"github.com/nicholasjackson/event-sauce/global"
	"github.com/nicholasjackson/event-sauce/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var queue *DeadLetterQueue
var mockDal *mocks.MockDal
var duration string
var deadLetters []*entities.DeadLetterItem

func getLetters() []*entities.DeadLetterItem {
	return deadLetters
}

func SetupDeadTests(t *testing.T) {

	mockDal = &mocks.MockDal{}
	queue = &DeadLetterQueue{Dal: mockDal}
	duration = "1d"

	global.Config.RetryIntervals = []string{duration}
	mockDal.Mock.On("UpsertDeadLetterItem", mock.Anything).Return(nil)
	mockDal.Mock.On("GetDeadLetterItemsReadyForRetry").Return(getLetters, nil)
}

func TestAddEventCreatesAValidDeadLetter(t *testing.T) {
	SetupDeadTests(t)

	event := &entities.Event{EventName: "soemthing"}
	d, _ := time.ParseDuration(duration)

	queue.AddEvent(event)

	assert.Equal(t, 1, mockDal.UpsertDeadLetter.FailureCount)
	assert.False(t, mockDal.UpsertDeadLetter.FirstFailureDate.IsZero())
	assert.Equal(t,
		mockDal.UpsertDeadLetter.FirstFailureDate.Add(d),
		mockDal.UpsertDeadLetter.NextRetryDate) // 1 day + first fail
	mockDal.Mock.AssertCalled(t, "UpsertDeadLetterItem", mock.Anything)
}

func TestStartConsumingGetsDeadLetters(t *testing.T) {
	SetupDeadTests(t)
	queue.runConsumer(1, func(item interface{}) {})

	mockDal.Mock.AssertCalled(t, "GetDeadLetterItemsReadyForRetry", mock.Anything)
}

func TestStartConsumingCallsCallbackForEachDeadLetter(t *testing.T) {
	SetupDeadTests(t)
	deadLetters = []*entities.DeadLetterItem{
		&entities.DeadLetterItem{},
		&entities.DeadLetterItem{},
		&entities.DeadLetterItem{},
	}

	callbackCount := 0

	queue.runConsumer(3, func(item interface{}) {
		callbackCount++
	})

	time.Sleep(20 * time.Millisecond)

	assert.Equal(t, 3, callbackCount)
}
