package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/facebookgo/inject"
	"github.com/nicholasjackson/event-sauce/global"
	"github.com/nicholasjackson/event-sauce/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type EventTestDependencies struct {
	StatsMock *mocks.MockStatsD `inject:"statsd"`
	QueueMock *mocks.MockQueue  `inject:"eventqueue"`
}

var mockEventDeps *EventTestDependencies

func SetupEventTest(t *testing.T) {
	EventHandlerDependencies = &EventDependencies{}
	mockEventDeps = &EventTestDependencies{}

	statsDMock := &mocks.MockStatsD{}
	queueMock := &mocks.MockQueue{}

	_ = global.SetupInjection(
		&inject.Object{Value: EventHandlerDependencies},
		&inject.Object{Value: mockEventDeps},
		&inject.Object{Value: statsDMock, Name: "statsd"},
		&inject.Object{Value: queueMock, Name: "eventqueue"},
	)

	mockEventDeps.StatsMock.Mock.On("Increment", mock.Anything).Return()
	mockEventDeps.QueueMock.Mock.On("Add", mock.Anything, mock.Anything).Return(nil)
}

func TestEventCallsStatsD(t *testing.T) {
	SetupEventTest(t)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(""))

	EventHandler(&responseRecorder, &request)

	mockEventDeps.StatsMock.Mock.AssertCalled(t, "Increment", EVENT_HANDLER_CALLED)
}

func TestEventWithNoPayloadReturns400(t *testing.T) {
	SetupEventTest(t)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(`
    {
      "message_name": "mymessage.stuff"
    }`))

	EventHandler(&responseRecorder, &request)

	assert.Equal(t, 400, responseRecorder.Code)
}

func TestEventWithNoMessageNameReturns400(t *testing.T) {
	SetupEventTest(t)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(`
    {
      "payload": {"name": "mymessage.stuff"}
    }`))

	EventHandler(&responseRecorder, &request)

	assert.Equal(t, 400, responseRecorder.Code)
}

func TestEventWithValidMessageReturns200(t *testing.T) {
	SetupEventTest(t)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(`
    {
      "message_name": "mymessage.stuff",
      "payload": {"name": "mymessage.stuff"}
    }`))

	EventHandler(&responseRecorder, &request)

	assert.Equal(t, 200, responseRecorder.Code)
}

func TestEventWithValidMessageAddsToQueue(t *testing.T) {
	SetupEventTest(t)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(`
    {
      "message_name": "mymessage.stuff",
      "payload": {"name": "mymessage.stuff"}
    }`))

	EventHandler(&responseRecorder, &request)

	mockEventDeps.QueueMock.Mock.AssertCalled(t, "Add", "mymessage.stuff", `{"name": "mymessage.stuff"}`)
}
