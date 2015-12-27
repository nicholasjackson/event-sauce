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

type RegisterTestDependencies struct {
	StatsMock *mocks.MockStatsD `inject:"statsd"`
}

var mockRegisterDeps *RegisterTestDependencies

func SetupRegisterTest(t *testing.T) {
	RegisterHandlerDependencies = &RegisterDependencies{}
	mockRegisterDeps = &RegisterTestDependencies{}

	statsDMock := &mocks.MockStatsD{}

	_ = global.SetupInjection(
		&inject.Object{Value: RegisterHandlerDependencies},
		&inject.Object{Value: mockRegisterDeps},
		&inject.Object{Value: statsDMock, Name: "statsd"},
	)

	mockRegisterDeps.StatsMock.Mock.On("Increment", mock.Anything).Return()
}

func TestCallsStatsD(t *testing.T) {
	SetupRegisterTest(t)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(""))

	RegisterHandler(&responseRecorder, &request)

	mockRegisterDeps.StatsMock.Mock.AssertCalled(t, "Increment", REGISTER_HANDLER_CALLED)
}

func TestRegisterReturnsBadRequest(t *testing.T) {
	SetupRegisterTest(t)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(""))

	RegisterHandler(&responseRecorder, &request)

	assert.Equal(t, responseRecorder.Code, 400)
}
