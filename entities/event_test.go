package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatesValidEventStoreItem(t *testing.T) {
	eventName := "mytest.event"
	payload := "mypayload"
	event := Event{EventName: eventName, Payload: payload}
	eventStoreItem := NewEventStoreItem(event)

	assert.Equal(t, eventName, eventStoreItem.Event.EventName)
	assert.Equal(t, payload, eventStoreItem.Event.Payload)
	assert.NotEmpty(t, eventStoreItem.Id)
	assert.NotEmpty(t, eventStoreItem.CreationDate)
}

func TestCreatesValidDeadLetterItem(t *testing.T) {
	eventName := "mytest.event"
	payload := "mypayload"
	event := Event{EventName: eventName, Payload: payload}
	deadLetterItem := NewDeadLetterItem(event)

	assert.Equal(t, eventName, deadLetterItem.Event.EventName)
	assert.Equal(t, payload, deadLetterItem.Event.Payload)
	assert.NotEmpty(t, deadLetterItem.Id)
	assert.NotEmpty(t, deadLetterItem.CreationDate)
	assert.True(t, deadLetterItem.FirstFailureDate.IsZero())
	assert.True(t, deadLetterItem.NextRetryDate.IsZero())
	assert.Equal(t, 0, deadLetterItem.FailureCount)
	assert.Empty(t, deadLetterItem.CallbackUrl)
}
