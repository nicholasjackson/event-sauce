package handlers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/facebookgo/inject"
	"github.com/nicholasjackson/event-sauce/global"
	"github.com/nicholasjackson/event-sauce/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type HealthTestDependencies struct {
	StatsMock *mocks.MockStatsD `inject:"statsd"`
}

var mockHealthDeps *HealthTestDependencies

func SetupHealthTest(t *testing.T) {
	mockHealthDeps = &HealthTestDependencies{}
	HealthHandlerDependencies = &HealthDependencies{}

	statsDMock := &mocks.MockStatsD{}

	_ = global.SetupInjection(
		&inject.Object{Value: HealthHandlerDependencies},
		&inject.Object{Value: mockHealthDeps},
		&inject.Object{Value: log.New(os.Stdout, "tester", log.Lshortfile)},
		&inject.Object{Value: statsDMock, Name: "statsd"},
	)

	mockHealthDeps.StatsMock.Mock.On("Increment", mock.Anything).Return()
}

// Simple test to show how we can use the ResponseRecorder to test our HTTP handlers
func TestHealthHandler(t *testing.T) {
	SetupHealthTest(t)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request

	HealthHandler(&responseRecorder, &request)

	assert.Equal(t, 200, responseRecorder.Code)
}

func TestHealthHandlerSetStats(t *testing.T) {
	SetupHealthTest(t)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request

	HealthHandler(&responseRecorder, &request)

	mockHealthDeps.StatsMock.Mock.AssertCalled(t, "Increment", HEALTH_HANDLER+GET+CALLED)
}
