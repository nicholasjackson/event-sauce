package main

import (
	"sync"
	"testing"

	"github.com/facebookgo/inject"
	"github.com/nicholasjackson/event-sauce/entities"
	"github.com/nicholasjackson/event-sauce/global"
	"github.com/nicholasjackson/event-sauce/mocks"
	"github.com/nicholasjackson/event-sauce/workers"
	"github.com/stretchr/testify/mock"
)

type MockWorkerFactory struct {
	mock.Mock
}

func (m *MockWorkerFactory) Create() workers.Worker {
	args := m.Mock.Called()
	return args.Get(0).(workers.Worker)
}

type ClientTestDependencies struct {
	StatsMock         *mocks.MockStatsD  `inject:"statsd"`
	DalMock           *mocks.MockDal     `inject:"dal"`
	EventQueueMock    *mocks.MockQueue   `inject:"eventqueue"`
	WorkerFactoryMock *MockWorkerFactory `inject:"eventqueueworkerfactory"`
}

var mockClientDeps *ClientTestDependencies
var mockWorker *mocks.MockWorker
var testWaitGroup sync.WaitGroup

func SetupClientTest(t *testing.T) {
	ClientDeps = &ClientDependencies{}
	mockClientDeps = &ClientTestDependencies{}

	statsDMock := &mocks.MockStatsD{}
	dalMock := &mocks.MockDal{}
	eventQueueMock := &mocks.MockQueue{}
	mockWorkerFactory := &MockWorkerFactory{}
	mockWorker = &mocks.MockWorker{}

	testWaitGroup = sync.WaitGroup{}
	testWaitGroup.Add(1)

	_ = global.SetupInjection(
		&inject.Object{Value: ClientDeps},
		&inject.Object{Value: mockClientDeps},
		&inject.Object{Value: statsDMock, Name: "statsd"},
		&inject.Object{Value: dalMock, Name: "dal"},
		&inject.Object{Value: eventQueueMock, Name: "eventqueue"},
		&inject.Object{Value: mockWorkerFactory, Name: "eventqueueworkerfactory"},
	)

	mockClientDeps.StatsMock.Mock.On("Increment", mock.Anything).Return()
	mockClientDeps.EventQueueMock.Mock.On("StartConsuming", mock.Anything, mock.Anything, mock.Anything)
	mockClientDeps.WorkerFactoryMock.Mock.On("Create").Return(mockWorker)
	mockWorker.Mock.On("HandleMessage", mock.Anything, mock.Anything).Return(nil)
	//mockRegisterDeps.DalMock.Mock.On("UpsertRegistration", mock.Anything).Return(nil, nil)
	//mockRegisterDeps.DalMock.Mock.On("DeleteRegistration", mock.Anything).Return(nil, nil)
}

func TestClientCreateCallsStatsD(t *testing.T) {
	SetupClientTest(t)

	startClient(&testWaitGroup)

	mockClientDeps.StatsMock.Mock.AssertCalled(t, "Increment", CLIENT_STARTED)
}

func TestClientStartsPolling(t *testing.T) {
	SetupClientTest(t)

	startClient(&testWaitGroup)

	mockClientDeps.EventQueueMock.Mock.AssertCalled(t, "StartConsuming", mock.Anything, mock.Anything, mock.Anything)
}

func TestClientCreatesWorkerWhenItemDeQueued(t *testing.T) {
	SetupClientTest(t)

	startClient(&testWaitGroup)
	mockClientDeps.EventQueueMock.ConsumerCallback(&entities.Event{})
	mockClientDeps.WorkerFactoryMock.Mock.AssertCalled(t, "Create")
}

func TestClientProcessesEventWhenItemDeQueued(t *testing.T) {
	SetupClientTest(t)

	startClient(&testWaitGroup)
	mockClientDeps.EventQueueMock.ConsumerCallback(&entities.Event{})
	mockWorker.Mock.AssertCalled(t, "HandleMessage", mock.Anything)
}
