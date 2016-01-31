package main

import (
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/facebookgo/inject"
	"github.com/nicholasjackson/sorcery/entities"
	"github.com/nicholasjackson/sorcery/global"
	"github.com/nicholasjackson/sorcery/handlers"
	"github.com/nicholasjackson/sorcery/mocks"
	"github.com/nicholasjackson/sorcery/workers"
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
	StatsMock                   *mocks.MockStatsD  `inject:"statsd"`
	EventQueueMock              *mocks.MockQueue   `inject:"eventqueue"`
	DeadLetterQueueMock         *mocks.MockQueue   `inject:"deadletterqueue"`
	EventWorkerFactoryMock      *MockWorkerFactory `inject:"eventqueueworkerfactory"`
	DeadLetterWorkerFactoryMock *MockWorkerFactory `inject:"deadletterqueueworkerfactory"`
}

var mockClientDeps *ClientTestDependencies
var mockWorker *mocks.MockWorker
var mockDeadLetterWorker *mocks.MockWorker
var testWaitGroup sync.WaitGroup

func SetupClientTest(t *testing.T) {
	ClientDeps = &ClientDependencies{}
	mockClientDeps = &ClientTestDependencies{}

	statsDMock := &mocks.MockStatsD{}
	eventQueueMock := &mocks.MockQueue{}
	deadLetterQueueMock := &mocks.MockQueue{}
	mockEventWorkerFactory := &MockWorkerFactory{}
	mockDeadLetterWorkerFactory := &MockWorkerFactory{}
	mockWorker = &mocks.MockWorker{}
	mockDeadLetterWorker = &mocks.MockWorker{}

	global.Config.RetryIntervals = []string{"10s"}

	testWaitGroup = sync.WaitGroup{}
	testWaitGroup.Add(1)

	_ = global.SetupInjection(
		&inject.Object{Value: ClientDeps},
		&inject.Object{Value: mockClientDeps},
		&inject.Object{Value: statsDMock, Name: "statsd"},
		&inject.Object{Value: log.New(os.Stdout, "Testing: ", log.Lshortfile)},
		&inject.Object{Value: eventQueueMock, Name: "eventqueue"},
		&inject.Object{Value: deadLetterQueueMock, Name: "deadletterqueue"},
		&inject.Object{Value: mockEventWorkerFactory, Name: "eventqueueworkerfactory"},
		&inject.Object{Value: mockDeadLetterWorkerFactory, Name: "deadletterqueueworkerfactory"},
	)

	mockClientDeps.StatsMock.Mock.On("Increment", mock.Anything).Return()
	mockClientDeps.EventQueueMock.Mock.On("StartConsuming", mock.Anything, mock.Anything, mock.Anything)
	mockClientDeps.EventWorkerFactoryMock.Mock.On("Create").Return(mockWorker)
	mockWorker.Mock.On("HandleItem", mock.Anything, mock.Anything).Return(nil)

	mockClientDeps.DeadLetterQueueMock.Mock.On("StartConsuming", mock.Anything, mock.Anything, mock.Anything)
	mockClientDeps.DeadLetterWorkerFactoryMock.Mock.On("Create").Return(mockDeadLetterWorker)
	mockDeadLetterWorker.Mock.On("HandleItem", mock.Anything, mock.Anything).Return(nil)
}

func TestEventQueueClientCreateCallsStatsD(t *testing.T) {
	SetupClientTest(t)

	startClient(&testWaitGroup)
	time.Sleep(10 * time.Millisecond) // wait for prcessEventQueue to start

	mockClientDeps.StatsMock.Mock.AssertCalled(t, "Increment", handlers.EVENT_QUEUE+handlers.STARTED)
}

func TestEventQueueClientStartsPolling(t *testing.T) {
	SetupClientTest(t)

	startClient(&testWaitGroup)
	time.Sleep(10 * time.Millisecond) // wait for prcessEventQueue to start

	mockClientDeps.EventQueueMock.Mock.AssertCalled(t, "StartConsuming", mock.Anything, mock.Anything, mock.Anything)
}

func TestEventQueueClientCreatesWorkerWhenItemDeQueued(t *testing.T) {
	SetupClientTest(t)

	startClient(&testWaitGroup)
	time.Sleep(10 * time.Millisecond) // wait for prcessEventQueue to start

	mockClientDeps.EventQueueMock.ConsumerCallback(&entities.Event{})
	mockClientDeps.EventWorkerFactoryMock.Mock.AssertCalled(t, "Create")
}

func TestEventQueueClientProcessesEventWhenItemDeQueued(t *testing.T) {
	SetupClientTest(t)

	startClient(&testWaitGroup)
	time.Sleep(10 * time.Millisecond) // wait for prcessEventQueue to start

	mockClientDeps.EventQueueMock.ConsumerCallback(&entities.Event{})
	mockWorker.Mock.AssertCalled(t, "HandleItem", mock.Anything)
}

func TestDeadLetterQueueClientCreateCallsStatsD(t *testing.T) {
	SetupClientTest(t)

	startClient(&testWaitGroup)
	time.Sleep(10 * time.Millisecond) // wait for prcessEventQueue to start

	mockClientDeps.StatsMock.Mock.AssertCalled(t, "Increment", handlers.DEAD_LETTER_QUEUE+handlers.STARTED)
}

func TestDeadLetterQueueClientStartsPolling(t *testing.T) {
	duration, _ := time.ParseDuration("10s")
	SetupClientTest(t)

	startClient(&testWaitGroup)
	time.Sleep(10 * time.Millisecond) // wait for prcessEventQueue to start

	mockClientDeps.DeadLetterQueueMock.Mock.AssertCalled(t, "StartConsuming", mock.Anything, duration, mock.Anything)
}

func TestDeadLetterQueueClientCreatesWorkerWhenItemDeQueued(t *testing.T) {
	SetupClientTest(t)

	startClient(&testWaitGroup)
	time.Sleep(10 * time.Millisecond) // wait for prcessEventQueue to start

	mockClientDeps.DeadLetterQueueMock.ConsumerCallback(&entities.Event{})
	mockClientDeps.DeadLetterWorkerFactoryMock.Mock.AssertCalled(t, "Create")
}

func TestDeadLetterQueueClientProcessesEventWhenItemDeQueued(t *testing.T) {
	SetupClientTest(t)

	startClient(&testWaitGroup)
	time.Sleep(10 * time.Millisecond) // wait for prcessEventQueue to start

	mockClientDeps.DeadLetterQueueMock.ConsumerCallback(&entities.Event{})
	mockDeadLetterWorker.Mock.AssertCalled(t, "HandleItem", mock.Anything)
}
