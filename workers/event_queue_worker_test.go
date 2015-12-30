package workers

import (
	"fmt"
	"testing"

	"github.com/nicholasjackson/event-sauce/entities"
	"github.com/nicholasjackson/event-sauce/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockDispatcher *mocks.MockEventDispatcher
var mockDal *mocks.MockDal
var mockQueue *mocks.MockQueue
var worker *EventQueueWorker

func setupTests(t *testing.T) {
	mockDispatcher = &mocks.MockEventDispatcher{}
	mockDal = &mocks.MockDal{}
	mockQueue = &mocks.MockQueue{}
	worker = New(mockDispatcher, mockDal, mockQueue)

	mockDispatcher.Mock.On("DispatchMessage", mock.Anything, mock.Anything).Return(200, nil)
	mockDal.Mock.On("GetRegistrationsByMessage", mock.Anything).Return(nil, nil)
	mockDal.Mock.On("DeleteRegistration", mock.Anything).Return(nil)
	mockDal.Mock.On("UpsertEvent", mock.Anything).Return(nil)
	mockQueue.Mock.On("AddEvent", mock.Anything).Return(nil)
}

func TestSetsMessageDispatcherAndDal(t *testing.T) {
	setupTests(t)

	assert.Equal(t, mockDispatcher, worker.eventDispatcher)
	assert.Equal(t, mockDal, worker.dal)
	assert.Equal(t, mockQueue, worker.deadLetterQueue)
}

func TestHandleMessageSavesToMessageStore(t *testing.T) {
	setupTests(t)
	message := &entities.Event{}

	worker.HandleMessage(message)

	mockDal.Mock.AssertCalled(t, "UpsertEvent", message)
}

func TestHandleMessageAttemptsToDispatchMessage(t *testing.T) {
	setupTests(t)
	message := &entities.Event{}
	endpoint := "myendpoint"
	registrations := []*entities.Registration{&entities.Registration{CallbackUrl: endpoint}}
	mockDal.Mock.ExpectedCalls = []*mock.Call{} // reset calls
	mockDal.Mock.On("GetRegistrationsByMessage", mock.Anything).Return(registrations, nil)
	mockDal.Mock.On("DeleteRegistration", mock.Anything).Return(nil)
	mockDal.Mock.On("UpsertEvent", mock.Anything).Return(nil)

	worker.HandleMessage(message)

	mockDispatcher.Mock.AssertCalled(t, "DispatchMessage", message, endpoint)
}

func TestHandleMessageAttemptsToDispatchMultipleMessage(t *testing.T) {
	setupTests(t)
	message := &entities.Event{}
	endpoint := "myendpoint"
	registrations := []*entities.Registration{
		&entities.Registration{CallbackUrl: endpoint},
		&entities.Registration{CallbackUrl: endpoint},
		&entities.Registration{CallbackUrl: endpoint}}
	mockDal.Mock.ExpectedCalls = []*mock.Call{} // reset calls
	mockDal.Mock.On("GetRegistrationsByMessage", mock.Anything).Return(registrations, nil)
	mockDal.Mock.On("DeleteRegistration", mock.Anything).Return(nil)
	mockDal.Mock.On("UpsertEvent", mock.Anything).Return(nil)

	worker.HandleMessage(message)

	mockDispatcher.Mock.AssertNumberOfCalls(t, "DispatchMessage", 3)
}

func TestHandleMessageGetsRegisterdEndpointsFromDB(t *testing.T) {
	setupTests(t)
	message := &entities.Event{MessageName: "mymessage"}

	worker.HandleMessage(message)

	mockDal.Mock.AssertCalled(t, "GetRegistrationsByMessage", "mymessage")
}

func TestDispatchMessageFailureRemovesRegistrationWhenRegistrationFoundAndEndpointDoesNotExist(t *testing.T) {
	setupTests(t)
	message := &entities.Event{}
	endpoint := "myendpoint"
	registrations := []*entities.Registration{&entities.Registration{CallbackUrl: endpoint}}

	mockDispatcher.Mock.ExpectedCalls = []*mock.Call{} // reset calls
	mockDispatcher.Mock.On("DispatchMessage", message, endpoint).Return(404, fmt.Errorf("Unable to complete"))

	mockDal.Mock.ExpectedCalls = []*mock.Call{} // reset calls
	mockDal.Mock.On("GetRegistrationsByMessage", mock.Anything).Return(registrations, nil)
	mockDal.Mock.On("DeleteRegistration", mock.Anything).Return(nil)
	mockDal.Mock.On("UpsertEvent", mock.Anything).Return(nil)

	_ = worker.HandleMessage(message)

	mockDal.Mock.AssertCalled(t, "DeleteRegistration", registrations[0])
}

func TestDispatchMessageFailureAddsMessageToDeadLetterQueueWhenEndpointInErrorState(t *testing.T) {
	setupTests(t)
	message := &entities.Event{}
	endpoint := "myendpoint"
	registrations := []*entities.Registration{&entities.Registration{CallbackUrl: endpoint}}

	mockDispatcher.Mock.ExpectedCalls = []*mock.Call{} // reset calls
	mockDispatcher.Mock.On("DispatchMessage", message, endpoint).Return(500, fmt.Errorf("Unable to complete"))

	mockDal.Mock.ExpectedCalls = []*mock.Call{} // reset calls
	mockDal.Mock.On("GetRegistrationsByMessage", mock.Anything).Return(registrations, nil)
	mockDal.Mock.On("DeleteRegistration", mock.Anything).Return(nil)
	mockDal.Mock.On("UpsertEvent", mock.Anything).Return(nil)

	_ = worker.HandleMessage(message)

	assert.Equal(t, endpoint, message.Callback)
	mockQueue.Mock.AssertCalled(t, "AddEvent", message)
}

func TestDispatchMessageOKDoesNothing(t *testing.T) {
	setupTests(t)
	message := &entities.Event{}
	endpoint := "myendpoint"
	registrations := []*entities.Registration{&entities.Registration{CallbackUrl: endpoint}}

	mockDispatcher.Mock.ExpectedCalls = []*mock.Call{} // reset calls
	mockDispatcher.Mock.On("DispatchMessage", message, endpoint).Return(200, fmt.Errorf("Unable to complete"))

	mockDal.Mock.ExpectedCalls = []*mock.Call{} // reset calls
	mockDal.Mock.On("GetRegistrationsByMessage", mock.Anything).Return(registrations, nil)
	mockDal.Mock.On("DeleteRegistration", mock.Anything).Return(nil)
	mockDal.Mock.On("UpsertEvent", mock.Anything).Return(nil)

	_ = worker.HandleMessage(message)

	mockDal.Mock.AssertNumberOfCalls(t, "DeleteRegistration", 0)
	mockQueue.Mock.AssertNumberOfCalls(t, "AddEvent", 0)
}
