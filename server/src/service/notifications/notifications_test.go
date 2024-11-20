package service

import (
	"fmt"
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

func (m *MockUsersConnector) GetUserNameAndImage(id string, header string) (string, string, error) {
	//TODO implement me
	panic("implement me")
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

func TestSendMilestoneNotificationHasTheRightSideEffects(t *testing.T) {
	mockDevicesDB := new(MockDevicesDatabase)
	mockUsersConnector := new(MockUsersConnector)
	mockFirebaseConnector := new(MockFirebaseConnector)
	notificationsService := NewNotificationService(mockDevicesDB, mockUsersConnector, mockFirebaseConnector)

	mockDevicesDB.On(
		"GetDevicesTokens",
		"userId",
	).Return([]string{"token1", "token2"}, nil)

	mockUsersConnector.On(
		"CheckUserExists",
		"userId",
		"Bearer token",
	).Return(true, nil)

	mockFirebaseConnector.On(
		"SendNotificationToUserDevices",
		[]string{"token1", "token2"},
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		map[string]string{"deeplink": "twitSnap://profile_profile?userId=" + "followerId"},
	).Return(nil)

	err := notificationsService.SendFollowerMilestoneNotification("userId", "followerId", "Bearer token")
	assert.Nil(t, err)
	mockUsersConnector.AssertExpectations(t)
	mockDevicesDB.AssertExpectations(t)
	mockFirebaseConnector.AssertExpectations(t)
}

func TestSendMessageNotification_GetDevicesTokensError(t *testing.T) {
	mockDevicesDB := new(MockDevicesDatabase)
	mockUsersConnector := new(MockUsersConnector)
	mockFirebaseConnector := new(MockFirebaseConnector)
	notificationsService := NewNotificationService(mockDevicesDB, mockUsersConnector, mockFirebaseConnector)

	mockDevicesDB.On(
		"GetDevicesTokens",
		"receiverId",
	).Return([]string{}, fmt.Errorf("database error"))

	err := notificationsService.SendNewMessageNotification("receiverId", "senderId", "body", "chatReference")
	assert.NotNil(t, err)
	assert.Equal(t, "error getting devices tokens: database error", err.Detail)
	mockDevicesDB.AssertExpectations(t)
}

func TestSendMessageNotification_SendNotificationError(t *testing.T) {
	mockDevicesDB := new(MockDevicesDatabase)
	mockUsersConnector := new(MockUsersConnector)
	mockFirebaseConnector := new(MockFirebaseConnector)
	notificationsService := NewNotificationService(mockDevicesDB, mockUsersConnector, mockFirebaseConnector)

	mockDevicesDB.On(
		"GetDevicesTokens",
		"receiverId",
	).Return([]string{"token1", "token2"}, nil)

	mockFirebaseConnector.On(
		"SendNotificationToUserDevices",
		[]string{"token1", "token2"},
		"New message",
		"body",
		map[string]string{"deeplink": "twitSnap://messages_chat?userId=senderId?refId=chatReference"},
	).Return(modelErrors.InternalServerError("firebase error"))

	err := notificationsService.SendNewMessageNotification("receiverId", "senderId", "body", "chatReference")
	assert.NotNil(t, err)
	assert.Equal(t, "error sending notification: firebase error", err.Detail)
	mockFirebaseConnector.AssertExpectations(t)
	mockDevicesDB.AssertExpectations(t)
}

func TestSendMilestoneNotification_UserDoesNotExist(t *testing.T) {
	mockDevicesDB := new(MockDevicesDatabase)
	mockUsersConnector := new(MockUsersConnector)
	mockFirebaseConnector := new(MockFirebaseConnector)
	notificationsService := NewNotificationService(mockDevicesDB, mockUsersConnector, mockFirebaseConnector)

	mockUsersConnector.On(
		"CheckUserExists",
		"userId",
		"Bearer token",
	).Return(false, nil)

	err := notificationsService.SendFollowerMilestoneNotification("userId", "10", "Bearer token")
	assert.NotNil(t, err)
	assert.Equal(t, "receiver does not exist", err.Detail)
	mockUsersConnector.AssertExpectations(t)
}

func TestSendMilestoneNotification_CheckUserExistsError(t *testing.T) {
	mockDevicesDB := new(MockDevicesDatabase)
	mockUsersConnector := new(MockUsersConnector)
	mockFirebaseConnector := new(MockFirebaseConnector)
	notificationsService := NewNotificationService(mockDevicesDB, mockUsersConnector, mockFirebaseConnector)

	mockUsersConnector.On(
		"CheckUserExists",
		"userId",
		"Bearer token",
	).Return(false, modelErrors.ExternalServiceError("external service error"))

	err := notificationsService.SendFollowerMilestoneNotification("userId", "10", "Bearer token")
	assert.NotNil(t, err)
	assert.Equal(t, "error checking user existence: external service error", err.Detail)
	mockUsersConnector.AssertExpectations(t)
}

func TestSendMilestoneNotification_GetDevicesTokensError(t *testing.T) {
	mockDevicesDB := new(MockDevicesDatabase)
	mockUsersConnector := new(MockUsersConnector)
	mockFirebaseConnector := new(MockFirebaseConnector)
	notificationsService := NewNotificationService(mockDevicesDB, mockUsersConnector, mockFirebaseConnector)

	mockUsersConnector.On(
		"CheckUserExists",
		"userId",
		"Bearer token",
	).Return(true, nil)

	mockDevicesDB.On(
		"GetDevicesTokens",
		"userId",
	).Return([]string{}, modelErrors.InternalServerError("database error"))

	err := notificationsService.SendFollowerMilestoneNotification("userId", "10", "Bearer token")
	assert.NotNil(t, err)
	assert.Equal(t, "error getting devices tokens: database error", err.Detail)
	mockUsersConnector.AssertExpectations(t)
	mockDevicesDB.AssertExpectations(t)
}

func TestSendMilestoneNotification_SendNotificationError(t *testing.T) {
	mockDevicesDB := new(MockDevicesDatabase)
	mockUsersConnector := new(MockUsersConnector)
	mockFirebaseConnector := new(MockFirebaseConnector)
	notificationsService := NewNotificationService(mockDevicesDB, mockUsersConnector, mockFirebaseConnector)

	mockUsersConnector.On(
		"CheckUserExists",
		"userId",
		"Bearer token",
	).Return(true, nil)

	mockDevicesDB.On(
		"GetDevicesTokens",
		"userId",
	).Return([]string{"token1", "token2"}, nil)

	mockFirebaseConnector.On(
		"SendNotificationToUserDevices",
		[]string{"token1", "token2"},
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		map[string]string{"deeplink": "twitSnap://profile_profile?userId=followerId"},
	).Return(modelErrors.InternalServerError("firebase error"))

	err := notificationsService.SendFollowerMilestoneNotification("userId", "followerId", "Bearer token")
	assert.NotNil(t, err)
	assert.Equal(t, "error sending notification: firebase error", err.Detail)
	mockUsersConnector.AssertExpectations(t)
	mockDevicesDB.AssertExpectations(t)
	mockFirebaseConnector.AssertExpectations(t)
}

func TestSendMentionNotificationHasTheRightSideEffects(t *testing.T) {
	mockDevicesDB := new(MockDevicesDatabase)
	mockUsersConnector := new(MockUsersConnector)
	mockFirebaseConnector := new(MockFirebaseConnector)
	notificationsService := NewNotificationService(mockDevicesDB, mockUsersConnector, mockFirebaseConnector)

	mockDevicesDB.On(
		"GetDevicesTokens",
		"userId",
	).Return([]string{"token1", "token2"}, nil)

	mockUsersConnector.On(
		"CheckUserExists",
		"userId",
		"Bearer token",
	).Return(true, nil)

	mockFirebaseConnector.On(
		"SendNotificationToUserDevices",
		[]string{"token1", "token2"},
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		map[string]string{"deeplink": "twitSnap://home_twitSnap?twitSnapId=" + "postId"},
	).Return(nil)

	err := notificationsService.SendMentionNotification("userId", "followerId", "postId", "Bearer token")
	assert.Nil(t, err)
	mockUsersConnector.AssertExpectations(t)
	mockDevicesDB.AssertExpectations(t)
	mockFirebaseConnector.AssertExpectations(t)
}
