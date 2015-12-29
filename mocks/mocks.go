package mocks

import (
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
	UpsertObject *entities.Registration
	DeleteObject *entities.Registration
}

func (m *MockDal) GetRegistrationByMessageAndCallback(message string, callback_url string) (*entities.Registration, error) {
	args := m.Mock.Called(message, callback_url)
	if args.Get(0) != nil {
		return args.Get(0).(*entities.Registration), args.Error(1)
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
