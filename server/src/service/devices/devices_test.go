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

func (m *MockUsersConnector) CheckUserExists(id string, header string) (bool, error) {
	args := m.Called(id, header)
	return args.Bool(0), args.Error(1)
}

type MockDevicesDatabase struct {
	mock.Mock
}

func (m *MockDevicesDatabase) AddDevice(id string, token string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockDevicesDatabase) GetDevicesTokens(id string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func TestDeviceService_AddDeviceFailsWhenUserDoesntExist(t *testing.T) {
	userConnectorMock := new(MockUsersConnector)
	devicesDatabaseMock := new(MockDevicesDatabase)
	devicesService := NewDeviceService(userConnectorMock, devicesDatabaseMock)
	userConnectorMock.On("CheckUserExists", "userId", "authHeader").Return(false, nil)

	err := devicesService.AddDevice("userId", "deviceToken", "authHeader")

	assert.Equal(t, modelErrors.BadRequestError("User does not exist"), err)

}
