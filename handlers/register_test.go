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
	mockRegisterDeps.DalMock.Mock.On("DeleteRegistration", mock.Anything).Return(nil, nil)
}

func TestRegisterCreateCallsStatsD(t *testing.T) {
	SetupRegisterTest(t)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(""))

	RegisterCreateHandler(&responseRecorder, &request)

	mockRegisterDeps.StatsMock.Mock.AssertCalled(t, "Increment", REGISTER_CREATE_HANDLER_CALLED)
}

func TestRegisterCreateWithNoPayloadReturnsBadRequest(t *testing.T) {
	SetupRegisterTest(t)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(""))

	RegisterCreateHandler(&responseRecorder, &request)

	assert.Equal(t, 400, responseRecorder.Code)
}

func TestRegisterCreateWithNoEventNameReturnsBadRequest(t *testing.T) {
	SetupRegisterTest(t)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(`
		{
			"callback_url": "dfdffd"
		}`))

	RegisterCreateHandler(&responseRecorder, &request)

	assert.Equal(t, 400, responseRecorder.Code)
}

func TestRegisterCreateWithNoCallbackUrlReturnsBadRequest(t *testing.T) {
	SetupRegisterTest(t)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(`
		{
			"event_name": "dfdffd"
		}`))

	RegisterCreateHandler(&responseRecorder, &request)

	assert.Equal(t, 400, responseRecorder.Code)
}

func TestRegisterCreateWithValidRequestSavesDataWhenRegistrationDoesNotExist(t *testing.T) {
	SetupRegisterTest(t)
	mockRegisterDeps.DalMock.Mock.On("GetRegistrationByEventAndCallback", "event.something", "http://some_callback_url.com").Return(nil, nil)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(`
		{
			"event_name": "event.something",
			"callback_url": "http://some_callback_url.com"
		}`))

	RegisterCreateHandler(&responseRecorder, &request)

	mockRegisterDeps.DalMock.Mock.AssertNumberOfCalls(t, "UpsertRegistration", 1)
	assert.Equal(t, 200, responseRecorder.Code)
}

func TestRegisterCreateWithValidRequestCreatesValidRegistration(t *testing.T) {
	SetupRegisterTest(t)
	mockRegisterDeps.DalMock.Mock.On("GetRegistrationByEventAndCallback", "event.something", "http://some_callback_url.com").Return(nil, nil)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(`
		{
			"event_name": "event.something",
			"callback_url": "http://some_callback_url.com"
		}`))

	RegisterCreateHandler(&responseRecorder, &request)

	registration := mockRegisterDeps.DalMock.UpsertObject
	assert.NotZero(t, registration.Id)
	assert.Equal(t, "event.something", registration.EventName)
	assert.Equal(t, "http://some_callback_url.com", registration.CallbackUrl)
}

func TestRegisterCreateWithValidRequestDoesNotSaveDataWhenRegistrationExists(t *testing.T) {
	SetupRegisterTest(t)
	mockRegisterDeps.DalMock.Mock.On("GetRegistrationByEventAndCallback", "event.something", "http://some_callback_url.com").Return(&entities.Registration{}, nil)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(`
		{
			"event_name": "event.something",
			"callback_url": "http://some_callback_url.com"
		}`))

	RegisterCreateHandler(&responseRecorder, &request)

	mockRegisterDeps.DalMock.Mock.AssertNumberOfCalls(t, "UpsertRegistration", 0)
	assert.Equal(t, 304, responseRecorder.Code)
}

func TestRegisterDeleteCallsStatsD(t *testing.T) {
	SetupRegisterTest(t)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(""))

	RegisterDeleteHandler(&responseRecorder, &request)

	mockRegisterDeps.StatsMock.Mock.AssertCalled(t, "Increment", REGISTER_DELETE_HANDLER_CALLED)
}

func TestRegisterDeleteWithNoPayloadReturnsBadRequest(t *testing.T) {
	SetupRegisterTest(t)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(""))

	RegisterDeleteHandler(&responseRecorder, &request)

	assert.Equal(t, 400, responseRecorder.Code)
}

func TestRegisterDeleteWithNoEventNameReturnsBadRequest(t *testing.T) {
	SetupRegisterTest(t)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(`
		{
			"callback_url": "dfdffd"
		}`))

	RegisterDeleteHandler(&responseRecorder, &request)

	assert.Equal(t, 400, responseRecorder.Code)
}

func TestRegisterDeleteWithNoCallbackUrlReturnsBadRequest(t *testing.T) {
	SetupRegisterTest(t)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(`
		{
			"event_name": "dfdffd"
		}`))

	RegisterDeleteHandler(&responseRecorder, &request)

	assert.Equal(t, 400, responseRecorder.Code)
}

func TestRegisterDeleteWithValidRequestReturns404WhenRegistrationDoesNotExist(t *testing.T) {
	SetupRegisterTest(t)
	mockRegisterDeps.DalMock.Mock.On(
		"GetRegistrationByEventAndCallback",
		"event.something",
		"http://some_callback_url.com").Return(nil, nil)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(`
		{
			"event_name": "event.something",
			"callback_url": "http://some_callback_url.com"
		}`))

	RegisterDeleteHandler(&responseRecorder, &request)

	mockRegisterDeps.DalMock.Mock.AssertNumberOfCalls(t, "DeleteRegistration", 0)
	assert.Equal(t, 404, responseRecorder.Code)
}

func TestRegisterDeleteWithValidRequestDeletesRegistration(t *testing.T) {
	SetupRegisterTest(t)
	registration := &entities.Registration{}
	mockRegisterDeps.DalMock.Mock.On(
		"GetRegistrationByEventAndCallback",
		"event.something",
		"http://some_callback_url.com").Return(registration, nil)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(`
		{
			"event_name": "event.something",
			"callback_url": "http://some_callback_url.com"
		}`))

	RegisterDeleteHandler(&responseRecorder, &request)

	assert.Equal(t, 200, responseRecorder.Code)
	mockRegisterDeps.DalMock.Mock.AssertCalled(t, "DeleteRegistration", registration)
}
