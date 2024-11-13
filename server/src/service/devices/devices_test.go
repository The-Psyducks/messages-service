package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	modelErrors "messages/src/model/errors"

	"testing"
)

type MockUsersConnector struct {
	mock.Mock
}

func (m *MockUsersConnector) GetUserNameAndImage(id string, header string) (string, string, error) {
	args := m.Called(id, header)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockUsersConnector) CheckUserExists(id string, header string) (bool, error) {
	args := m.Called(id, header)
	return args.Bool(0), args.Error(1)
}

type MockDevicesDatabase struct {
	mock.Mock
}

func (m *MockDevicesDatabase) AddDevice(id string, token string) error {
	args := m.Called(id, token)
	return args.Error(0)
}

func (m *MockDevicesDatabase) GetDevicesTokens(id string) ([]string, error) {
	args := m.Called(id)
	return args.Get(0).([]string), args.Error(1)
}

func TestDeviceService_AddDeviceFailsWhenUserDoesntExist(t *testing.T) {
	userConnectorMock := new(MockUsersConnector)
	devicesDatabaseMock := new(MockDevicesDatabase)
	devicesService := NewDeviceService(userConnectorMock, devicesDatabaseMock)
	userConnectorMock.On("CheckUserExists", "userId", "authHeader").Return(false, nil)

	err := devicesService.AddDevice("userId", "deviceToken", "authHeader")

	assert.Equal(t, modelErrors.BadRequestError("User does not exist"), err)

}

func TestDeviceService_AddDeviceFailsWhenUserValidationFails(t *testing.T) {
	userConnectorMock := new(MockUsersConnector)
	devicesDatabaseMock := new(MockDevicesDatabase)
	devicesService := NewDeviceService(userConnectorMock, devicesDatabaseMock)
	userConnectorMock.On("CheckUserExists", "userId", "authHeader").Return(false, modelErrors.ExternalServiceError("external service error"))

	err := devicesService.AddDevice("userId", "deviceToken", "authHeader")

	assert.Equal(t, modelErrors.ExternalServiceError("Error validating user: external service error"), err)
}

func TestDeviceService_AddDeviceFailsWhenDatabaseFails(t *testing.T) {
	userConnectorMock := new(MockUsersConnector)
	devicesDatabaseMock := new(MockDevicesDatabase)
	devicesService := NewDeviceService(userConnectorMock, devicesDatabaseMock)
	userConnectorMock.On("CheckUserExists", "userId", "authHeader").Return(true, nil)
	devicesDatabaseMock.On("AddDevice", "userId", "deviceToken").Return(modelErrors.InternalServerError("database error"))

	err := devicesService.AddDevice("userId", "deviceToken", "authHeader")

	assert.Equal(t, modelErrors.InternalServerError("Error adding device to db: database error"), err)
}

func TestDeviceService_AddDeviceHappyPath(t *testing.T) {
	userConnectorMock := new(MockUsersConnector)
	devicesDatabaseMock := new(MockDevicesDatabase)
	devicesService := NewDeviceService(userConnectorMock, devicesDatabaseMock)

	userConnectorMock.On("CheckUserExists", "userId", "authHeader").Return(true, nil)
	devicesDatabaseMock.On("AddDevice", "userId", "deviceToken").Return(nil)

	err := devicesService.AddDevice("userId", "deviceToken", "authHeader")

	assert.Nil(t, err)
}
