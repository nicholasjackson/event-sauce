package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/facebookgo/inject"
	"github.com/nicholasjackson/event-sauce/entities"
	"github.com/nicholasjackson/event-sauce/global"
	"github.com/nicholasjackson/event-sauce/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type RegisterTestDependencies struct {
	StatsMock *mocks.MockStatsD `inject:"statsd"`
	DalMock   *mocks.MockDal    `inject:"dal"`
}

var mockRegisterDeps *RegisterTestDependencies

func SetupRegisterTest(t *testing.T) {
	RegisterHandlerDependencies = &RegisterDependencies{}
	mockRegisterDeps = &RegisterTestDependencies{}

	statsDMock := &mocks.MockStatsD{}
	dalMock := &mocks.MockDal{}

	_ = global.SetupInjection(
		&inject.Object{Value: RegisterHandlerDependencies},
		&inject.Object{Value: mockRegisterDeps},
		&inject.Object{Value: statsDMock, Name: "statsd"},
		&inject.Object{Value: dalMock, Name: "dal"},
	)

	mockRegisterDeps.StatsMock.Mock.On("Increment", mock.Anything).Return()
	mockRegisterDeps.DalMock.Mock.On("UpsertRegistration", mock.Anything).Return(nil, nil)
}

func TestCallsStatsD(t *testing.T) {
	SetupRegisterTest(t)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(""))

	RegisterHandler(&responseRecorder, &request)

	mockRegisterDeps.StatsMock.Mock.AssertCalled(t, "Increment", REGISTER_HANDLER_CALLED)
}

func TestRegisterWithNoPayloadReturnsBadRequest(t *testing.T) {
	SetupRegisterTest(t)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(""))

	RegisterHandler(&responseRecorder, &request)

	assert.Equal(t, 400, responseRecorder.Code)
}

func TestRegisterWithNoMessageNameReturnsBadRequest(t *testing.T) {
	SetupRegisterTest(t)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(`
		{
			"callback_url": "dfdffd"
		}`))

	RegisterHandler(&responseRecorder, &request)

	assert.Equal(t, 400, responseRecorder.Code)
}

func TestRegisterWithNoCallbackUrlReturnsBadRequest(t *testing.T) {
	SetupRegisterTest(t)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(`
		{
			"message_name": "dfdffd"
		}`))

	RegisterHandler(&responseRecorder, &request)

	assert.Equal(t, 400, responseRecorder.Code)
}

func TestRegisterWithValidRequestSavesDataWhenRegistrationDoesNotExist(t *testing.T) {
	SetupRegisterTest(t)
	mockRegisterDeps.DalMock.Mock.On("GetRegistrationByMessageAndCallback", "event.something", "http://some_callback_url.com").Return(nil, nil)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(`
		{
			"message_name": "event.something",
			"callback_url": "http://some_callback_url.com"
		}`))

	RegisterHandler(&responseRecorder, &request)

	mockRegisterDeps.DalMock.Mock.AssertNumberOfCalls(t, "UpsertRegistration", 1)
	assert.Equal(t, 200, responseRecorder.Code)
}

func TestRegisterWithValidRequestCreatesValidRegistration(t *testing.T) {
	SetupRegisterTest(t)
	mockRegisterDeps.DalMock.Mock.On("GetRegistrationByMessageAndCallback", "event.something", "http://some_callback_url.com").Return(nil, nil)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(`
		{
			"message_name": "event.something",
			"callback_url": "http://some_callback_url.com"
		}`))

	RegisterHandler(&responseRecorder, &request)

	registration := mockRegisterDeps.DalMock.UpsertObject
	assert.Equal(t, "event.something", registration.MessageName)
	assert.Equal(t, "http://some_callback_url.com", registration.CallbackUrl)
}

func TestRegisterWithValidRequestDoesNotSaveDataWhenRegistrationExists(t *testing.T) {
	SetupRegisterTest(t)
	mockRegisterDeps.DalMock.Mock.On("GetRegistrationByMessageAndCallback", "event.something", "http://some_callback_url.com").Return(&entities.Registration{}, nil)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(`
		{
			"message_name": "event.something",
			"callback_url": "http://some_callback_url.com"
		}`))

	RegisterHandler(&responseRecorder, &request)

	mockRegisterDeps.DalMock.Mock.AssertNumberOfCalls(t, "UpsertRegistration", 0)
	assert.Equal(t, 304, responseRecorder.Code)
}
