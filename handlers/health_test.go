package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/facebookgo/inject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


type MockStatsD struct {
	mock.Mock
}

func (m *MockStatsD) Increment(label string) {
	_ = m.Mock.Called(label)
}

var statsDMock *MockStatsD


func TestSetup(t *testing.T) {
	// create an injection graph containing the mocked elements we wish to replace

	var g inject.Graph

	statsDMock = &MockStatsD{}

	err := g.Provide(
		&inject.Object{Value: HealthHandlerDependencies},
		&inject.Object{Value: statsDMock, Name: "statsd"},
	)

	if err != nil {
		fmt.Println(err)
	}

	if err := g.Populate(); err != nil {
		fmt.Println(err)
	}

	statsDMock.Mock.On("Increment", mock.Anything).Return()
}

// Simple test to show how we can use the ResponseRecorder to test our HTTP handlers
func TestHealthHandler(t *testing.T) {
	var responseRecorder httptest.ResponseRecorder
	var request http.Request

	HealthHandler(&responseRecorder, &request)

	assert.Equal(t, 200, responseRecorder.Code)
}


func TestHealthHandlerSetStats(t *testing.T) {
	var responseRecorder httptest.ResponseRecorder
	var request http.Request

	HealthHandler(&responseRecorder, &request)

	statsDMock.Mock.AssertCalled(t, "Increment", HEALTH_HANDLER_CALLED)
}

