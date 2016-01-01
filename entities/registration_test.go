package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatesCorrectly(t *testing.T) {
	registration := CreateNewRegistration("myevent", "mycallback")

	assert.NotZero(t, registration.Id)
	assert.NotZero(t, registration.CreationDate)
	assert.Equal(t, "myevent", registration.EventName)
	assert.Equal(t, "mycallback", registration.CallbackUrl)
}
