package mocks

import (
	"time"

	"github.com/nicholasjackson/event-sauce/entities"
	"github.com/stretchr/testify/mock"
)

type MockStatsD struct {
	mock.Mock
}

func (m *MockStatsD) Increment(label string) {
	_ = m.Mock.Called(label)
}

type MockDal struct {
	mock.Mock
	UpsertObject     *entities.Registration
	DeleteObject     *entities.Registration
	UpsertDeadLetter *entities.DeadLetterItem
}

func (m *MockDal) GetRegistrationsByEvent(event string) ([]*entities.Registration, error) {
	args := m.Mock.Called(event)
	if args.Get(0) != nil {
		f, ok := args.Get(0).(func() []*entities.Registration)
		if ok {
			return f(), args.Error(1)
		} else {
			return args.Get(0).([]*entities.Registration), args.Error(1)
		}
	} else {
		return nil, args.Error(1)
	}
}

func (m *MockDal) GetRegistrationByEventAndCallback(event string, callback_url string) (*entities.Registration, error) {
	args := m.Mock.Called(event, callback_url)
	if args.Get(0) != nil {
		f, ok := args.Get(0).(func() *entities.Registration)
		if ok {
			return f(), args.Error(1)
		} else {
			return args.Get(0).(*entities.Registration), args.Error(1)
		}
	} else {
		return nil, args.Error(1)
	}
}

func (m *MockDal) UpsertRegistration(registration *entities.Registration) error {
	args := m.Mock.Called(registration)
	m.UpsertObject = registration
	return args.Error(0)
}

func (m *MockDal) DeleteRegistration(registration *entities.Registration) error {
	args := m.Mock.Called(registration)
	m.DeleteObject = registration
	return args.Error(0)
}

func (m *MockDal) UpsertEventStore(event *entities.EventStoreItem) error {
	args := m.Mock.Called(event)

	return args.Error(0)
}

func (m *MockDal) UpsertDeadLetterItem(dead *entities.DeadLetterItem) error {
	m.UpsertDeadLetter = dead
	args := m.Mock.Called(dead)

	return args.Error(0)
}

func (m *MockDal) GetDeadLetterItemsReadyForRetry() ([]*entities.DeadLetterItem, error) {
	args := m.Mock.Called()

	if args.Get(0) != nil {
		f, ok := args.Get(0).(func() []*entities.DeadLetterItem)
		if ok {
			return f(), args.Error(1)
		} else {
			return args.Get(0).([]*entities.DeadLetterItem), args.Error(1)
		}

	} else {
		return nil, args.Error(1)
	}
}

func (m *MockDal) DeleteDeadLetterItems(dead []*entities.DeadLetterItem) error {
	args := m.Mock.Called(dead)

	return args.Error(0)
}

type MockQueue struct {
	mock.Mock
	ConsumerCallback func(callbackItem interface{})
}

func (m *MockQueue) Add(event_name string, payload string) error {
	args := m.Mock.Called(event_name, payload)
	return args.Error(0)
}

func (m *MockQueue) AddEvent(event *entities.Event, callback string) error {
	args := m.Mock.Called(event, callback)
	return args.Error(0)
}

func (m *MockQueue) StartConsuming(size int, poll_interval time.Duration, callback func(callbackItem interface{})) {
	m.ConsumerCallback = callback
	_ = m.Mock.Called(size, poll_interval, callback)
}

type MockWorker struct {
	mock.Mock
}

func (m *MockWorker) HandleItem(item interface{}) error {
	args := m.Mock.Called(item)
	return args.Error(0)
}

type MockEventDispatcher struct {
	mock.Mock
}

func (m *MockEventDispatcher) DispatchEvent(event *entities.Event, endpoint string) (int, error) {
	args := m.Mock.Called(event, endpoint)
	return args.Int(0), args.Error(1)
}
