package mocks

import "github.com/stretchr/testify/mock"

type MockStatsD struct {
	mock.Mock
}

func (m *MockStatsD) Increment(label string) {
	_ = m.Mock.Called(label)
}
