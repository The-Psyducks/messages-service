package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func (m *MockFirebaseConnector) SendNotificationToUserDevices(devicesTokens []string, title, body string, data map[string]string) error {
	args := m.Called(devicesTokens, title, body, data)
	return args.Error(0)
}

func TestSendMessageNotificationHasTheRightSideEffects(t *testing.T) {
	mockDevicesDB := new(MockDevicesDatabase)
	mockUsersConnector := new(MockUsersConnector)
	mockFirebaseConnector := new(MockFirebaseConnector)
	notificationsService := NewNotificationService(mockDevicesDB, mockUsersConnector, mockFirebaseConnector)

	mockDevicesDB.On(
		"GetDevicesTokens",
		"receiverId",
	).Return([]string{"token1", "token2"}, nil)

	//mockUsersConnector.On(
	//	"CheckUserExists",
	//	"receiverId",
	//	"Bearer token",
	//).Return(true, nil)

	mockFirebaseConnector.On(
		"SendNotificationToUserDevices",
		[]string{"token1", "token2"},
		"New message",
		"body",
		map[string]string{"deeplink": "twitSnap://messages_chat?userId=senderId?refId=chatReference"},
	).Return(nil)

	//err := notificationsService.SendNotificationTo("receiverId", "title", "body", "Bearer token")
	err := notificationsService.SendNewMessageNotification("receiverId", "senderId", "body", "chatReference")
	assert.Nil(t, err)
	mockUsersConnector.AssertExpectations(t)
	mockDevicesDB.AssertExpectations(t)
	mockFirebaseConnector.AssertExpectations(t)
}
