package service

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	modelErrors "messages/src/model/errors"
	"testing"
)

// Mock for DevicesDatabaseInterface
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

// Mock for UsersConnectorInterface
type MockUsersConnector struct {
	mock.Mock
}

func (m *MockUsersConnector) CheckUserExists(userID, authHeader string) (bool, error) {
	args := m.Called(userID, authHeader)
	return args.Bool(0), args.Error(1)
}

// Mock for FirebaseConnectorInterface
type MockFirebaseConnector struct {
	mock.Mock
}

func (m *MockFirebaseConnector) SendNotificationToUserDevices(tokens []string, title, body string) error {
	args := m.Called(tokens, title, body)
	return args.Error(0)
}

func TestSendNotification_HappyPath(t *testing.T) {
	mockDevicesDB := new(MockDevicesDatabase)
	mockUsersConnector := new(MockUsersConnector)
	mockFirebaseConnector := new(MockFirebaseConnector)
	service := NewNotificationService(mockDevicesDB, mockUsersConnector, mockFirebaseConnector)

	mockUsersConnector.On("CheckUserExists", "receiverId", "Bearer token").Return(true, nil)
	mockDevicesDB.On("GetDevicesTokens", "receiverId").Return([]string{"token1", "token2"}, nil)
	mockFirebaseConnector.On("SendNotificationToUserDevices", []string{"token1", "token2"}, "title", "body").Return(nil)

	err := service.SendNotification("receiverId", "title", "body", "Bearer token")

	assert.Nil(t, err)
	mockUsersConnector.AssertExpectations(t)
	mockDevicesDB.AssertExpectations(t)
	mockFirebaseConnector.AssertExpectations(t)
}

func TestSendNotification_UserDoesNotExist(t *testing.T) {
	mockDevicesDB := new(MockDevicesDatabase)
	mockUsersConnector := new(MockUsersConnector)
	mockFirebaseConnector := new(MockFirebaseConnector)
	service := NewNotificationService(mockDevicesDB, mockUsersConnector, mockFirebaseConnector)

	mockUsersConnector.On("CheckUserExists", "receiverId", "Bearer token").Return(false, nil)

	err := service.SendNotification("receiverId", "title", "body", "Bearer token")

	assert.NotNil(t, err)
	assert.Equal(t, "receiver does not exist", err.Detail)
	mockUsersConnector.AssertExpectations(t)
}

func TestSendNotification_ExternalServiceError(t *testing.T) {
	mockDevicesDB := new(MockDevicesDatabase)
	mockUsersConnector := new(MockUsersConnector)
	mockFirebaseConnector := new(MockFirebaseConnector)
	service := NewNotificationService(mockDevicesDB, mockUsersConnector, mockFirebaseConnector)

	mockUsersConnector.On("CheckUserExists", "receiverId", "Bearer token").Return(false, modelErrors.ExternalServiceError("external service error"))

	err := service.SendNotification("receiverId", "title", "body", "Bearer token")

	assert.NotNil(t, err)
	assert.Equal(t, "error checking user existence: external service error", err.Detail)
	mockUsersConnector.AssertExpectations(t)
}

func TestSendNotification_GetDevicesTokensError(t *testing.T) {
	mockDevicesDB := new(MockDevicesDatabase)
	mockUsersConnector := new(MockUsersConnector)
	mockFirebaseConnector := new(MockFirebaseConnector)
	service := NewNotificationService(mockDevicesDB, mockUsersConnector, mockFirebaseConnector)

	mockUsersConnector.On("CheckUserExists", "receiverId", "Bearer token").Return(true, nil)

	mockDevicesDB.On("GetDevicesTokens", "receiverId").Return([]string{}, errors.New("database error"))

	err := service.SendNotification("receiverId", "title", "body", "Bearer token")

	assert.NotNil(t, err)
	assert.Equal(t, "error getting devices tokens: database error", err.Detail)
	mockUsersConnector.AssertExpectations(t)
	mockDevicesDB.AssertExpectations(t)
}

func TestSendNotification_SendNotificationError(t *testing.T) {
	mockDevicesDB := new(MockDevicesDatabase)
	mockUsersConnector := new(MockUsersConnector)
	mockFirebaseConnector := new(MockFirebaseConnector)
	service := NewNotificationService(mockDevicesDB, mockUsersConnector, mockFirebaseConnector)

	mockUsersConnector.On("CheckUserExists", "receiverId", "Bearer token").Return(true, nil)
	mockDevicesDB.On("GetDevicesTokens", "receiverId").Return([]string{"token1", "token2"}, nil)
	mockFirebaseConnector.On("SendNotificationToUserDevices", []string{"token1", "token2"}, "title", "body").Return(modelErrors.InternalServerError("firebase error"))

	err := service.SendNotification("receiverId", "title", "body", "Bearer token")

	assert.NotNil(t, err)
	assert.Equal(t, "error sending notification: firebase error", err.Detail)
	mockUsersConnector.AssertExpectations(t)
	mockDevicesDB.AssertExpectations(t)
	mockFirebaseConnector.AssertExpectations(t)
}
