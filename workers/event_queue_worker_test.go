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

	mockDispatcher.Mock.On("DispatchEvent", mock.Anything, mock.Anything).Return(200, nil)
	mockDal.Mock.On("GetRegistrationsByEvent", mock.Anything).Return(nil, nil)
	mockDal.Mock.On("DeleteRegistration", mock.Anything).Return(nil)
	mockDal.Mock.On("UpsertEvent", mock.Anything).Return(nil)
	mockQueue.Mock.On("AddEvent", mock.Anything).Return(nil)
}

func TestSetsEventDispatcherAndDal(t *testing.T) {
	setupTests(t)

	assert.Equal(t, mockDispatcher, worker.eventDispatcher)
	assert.Equal(t, mockDal, worker.dal)
	assert.Equal(t, mockQueue, worker.deadLetterQueue)
}

func TestHandleEventSavesToEventStore(t *testing.T) {
	setupTests(t)
	event := &entities.Event{}

	worker.HandleEvent(event)

	mockDal.Mock.AssertCalled(t, "UpsertEvent", event)
}

func TestHandleEventAttemptsToDispatchEvent(t *testing.T) {
	setupTests(t)
	event := &entities.Event{}
	endpoint := "myendpoint"
	registrations := []*entities.Registration{&entities.Registration{CallbackUrl: endpoint}}
	mockDal.Mock.ExpectedCalls = []*mock.Call{} // reset calls
	mockDal.Mock.On("GetRegistrationsByEvent", mock.Anything).Return(registrations, nil)
	mockDal.Mock.On("DeleteRegistration", mock.Anything).Return(nil)
	mockDal.Mock.On("UpsertEvent", mock.Anything).Return(nil)

	worker.HandleEvent(event)

	mockDispatcher.Mock.AssertCalled(t, "DispatchEvent", event, endpoint)
}

func TestHandleEventAttemptsToDispatchMultipleEvent(t *testing.T) {
	setupTests(t)
	event := &entities.Event{}
	endpoint := "myendpoint"
	registrations := []*entities.Registration{
		&entities.Registration{CallbackUrl: endpoint},
		&entities.Registration{CallbackUrl: endpoint},
		&entities.Registration{CallbackUrl: endpoint}}
	mockDal.Mock.ExpectedCalls = []*mock.Call{} // reset calls
	mockDal.Mock.On("GetRegistrationsByEvent", mock.Anything).Return(registrations, nil)
	mockDal.Mock.On("DeleteRegistration", mock.Anything).Return(nil)
	mockDal.Mock.On("UpsertEvent", mock.Anything).Return(nil)

	worker.HandleEvent(event)

	mockDispatcher.Mock.AssertNumberOfCalls(t, "DispatchEvent", 3)
}

func TestHandleEventGetsRegisterdEndpointsFromDB(t *testing.T) {
	setupTests(t)
	event := &entities.Event{EventName: "myevent"}

	worker.HandleEvent(event)

	mockDal.Mock.AssertCalled(t, "GetRegistrationsByEvent", "myevent")
}

func TestDispatchEventFailureRemovesRegistrationWhenRegistrationFoundAndEndpointDoesNotExist(t *testing.T) {
	setupTests(t)
	event := &entities.Event{}
	endpoint := "myendpoint"
	registrations := []*entities.Registration{&entities.Registration{CallbackUrl: endpoint}}

	mockDispatcher.Mock.ExpectedCalls = []*mock.Call{} // reset calls
	mockDispatcher.Mock.On("DispatchEvent", event, endpoint).Return(404, fmt.Errorf("Unable to complete"))

	mockDal.Mock.ExpectedCalls = []*mock.Call{} // reset calls
	mockDal.Mock.On("GetRegistrationsByEvent", mock.Anything).Return(registrations, nil)
	mockDal.Mock.On("DeleteRegistration", mock.Anything).Return(nil)
	mockDal.Mock.On("UpsertEvent", mock.Anything).Return(nil)

	_ = worker.HandleEvent(event)

	mockDal.Mock.AssertCalled(t, "DeleteRegistration", registrations[0])
}

func TestDispatchEventFailureAddsEventToDeadLetterQueueWhenEndpointInErrorState(t *testing.T) {
	setupTests(t)
	event := &entities.Event{}
	endpoint := "myendpoint"
	registrations := []*entities.Registration{&entities.Registration{CallbackUrl: endpoint}}

	mockDispatcher.Mock.ExpectedCalls = []*mock.Call{} // reset calls
	mockDispatcher.Mock.On("DispatchEvent", event, endpoint).Return(500, fmt.Errorf("Unable to complete"))

	mockDal.Mock.ExpectedCalls = []*mock.Call{} // reset calls
	mockDal.Mock.On("GetRegistrationsByEvent", mock.Anything).Return(registrations, nil)
	mockDal.Mock.On("DeleteRegistration", mock.Anything).Return(nil)
	mockDal.Mock.On("UpsertEvent", mock.Anything).Return(nil)

	_ = worker.HandleEvent(event)

	assert.Equal(t, endpoint, event.Callback)
	mockQueue.Mock.AssertCalled(t, "AddEvent", event)
}

func TestDispatchEventOKDoesNothing(t *testing.T) {
	setupTests(t)
	event := &entities.Event{}
	endpoint := "myendpoint"
	registrations := []*entities.Registration{&entities.Registration{CallbackUrl: endpoint}}

	mockDispatcher.Mock.ExpectedCalls = []*mock.Call{} // reset calls
	mockDispatcher.Mock.On("DispatchEvent", event, endpoint).Return(200, fmt.Errorf("Unable to complete"))

	mockDal.Mock.ExpectedCalls = []*mock.Call{} // reset calls
	mockDal.Mock.On("GetRegistrationsByEvent", mock.Anything).Return(registrations, nil)
	mockDal.Mock.On("DeleteRegistration", mock.Anything).Return(nil)
	mockDal.Mock.On("UpsertEvent", mock.Anything).Return(nil)

	_ = worker.HandleEvent(event)

	mockDal.Mock.AssertNumberOfCalls(t, "DeleteRegistration", 0)
	mockQueue.Mock.AssertNumberOfCalls(t, "AddEvent", 0)
}
