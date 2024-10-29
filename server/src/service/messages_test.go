package service

import (
	goErrors "errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"messages/src/model/errors"
	"testing"
)

// Mock for RealTimeDatabaseInterface
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) SendMessage(senderId, receiverId, content string) (string, error) {
	args := m.Called(senderId, receiverId, content)
	return args.String(0), args.Error(1)
}

func (m *MockDatabase) GetConversations(id string) ([]string, error) {
	panic("implement me")
}

// Mock for ConnectorInterface
type MockUserConnector struct {
	mock.Mock
}

func (m *MockUserConnector) CheckUserExists(userID, authHeader string) (bool, error) {
	args := m.Called(userID, authHeader)
	return args.Bool(0), args.Error(1)
}

func TestSendMessage_HappyPath(t *testing.T) {
	mockDB := new(MockDatabase)
	mockUserConnector := new(MockUserConnector)
	expectedRef := "messageRef123"
	service := NewMessageService(mockDB, mockUserConnector)

	mockUserConnector.On("CheckUserExists", "existing_user_1", "Bearer token").Return(true, nil)
	mockUserConnector.On("CheckUserExists", "existing_user_2", "Bearer token").Return(true, nil)
	mockDB.On("SendMessage", "existing_user_1", "existing_user_2", "Hello, World!").Return(expectedRef, nil)

	ref, err := service.SendMessage("existing_user_1", "existing_user_2", "Hello, World!", "Bearer token")

	assert.Nil(t, err)
	assert.Equal(t, expectedRef, ref)

	mockDB.AssertExpectations(t)
	mockUserConnector.AssertExpectations(t)
}

func TestSendMessage_SenderValidationFails(t *testing.T) {
	mockDB := new(MockDatabase)
	mockUserConnector := new(MockUserConnector)
	service := NewMessageService(mockDB, mockUserConnector)

	mockUserConnector.On("CheckUserExists", "nonexistent_user", "Bearer token").Return(false, nil)

	ref, err := service.SendMessage("nonexistent_user", "existing_user_2", "Hello, World!", "Bearer token")

	assert.Equal(t, "", ref)
	assert.NotNil(t, err)
	assert.Equal(t, "sender does not exist", err.Detail)

	mockUserConnector.AssertExpectations(t)
}

func TestSendMessage_SenderExternalServiceError(t *testing.T) {
	mockDB := new(MockDatabase)
	mockUserConnector := new(MockUserConnector)
	service := NewMessageService(mockDB, mockUserConnector)

	mockUserConnector.On("CheckUserExists", "existing_user_1", "Bearer token").Return(false, errors.ExternalServiceError("external service error"))

	ref, err := service.SendMessage("existing_user_1", "existing_user_2", "Hello, World!", "Bearer token")

	assert.Equal(t, "", ref)
	assert.NotNil(t, err)
	assert.Equal(t, "error validating user: external service error", err.Detail)

	mockUserConnector.AssertExpectations(t)
}

func TestSendMessage_ReceiverValidationFails(t *testing.T) {
	mockDB := new(MockDatabase)
	mockUserConnector := new(MockUserConnector)
	service := NewMessageService(mockDB, mockUserConnector)

	mockUserConnector.On("CheckUserExists", "existing_user_1", "Bearer token").Return(true, nil)
	mockUserConnector.On("CheckUserExists", "nonexistent_user", "Bearer token").Return(false, nil)

	ref, err := service.SendMessage("existing_user_1", "nonexistent_user", "Hello, World!", "Bearer token")

	assert.Equal(t, "", ref)
	assert.NotNil(t, err)
	assert.Equal(t, "receiver does not exist", err.Detail)

	mockUserConnector.AssertExpectations(t)
}

func TestSendMessage_ReceiverExternalServiceError(t *testing.T) {
	mockDB := new(MockDatabase)
	mockUserConnector := new(MockUserConnector)
	service := NewMessageService(mockDB, mockUserConnector)

	mockUserConnector.On("CheckUserExists", "existing_user_1", "Bearer token").Return(true, nil)
	mockUserConnector.On("CheckUserExists", "existing_user_2", "Bearer token").Return(false, errors.ExternalServiceError("external service error"))

	ref, err := service.SendMessage("existing_user_1", "existing_user_2", "Hello, World!", "Bearer token")

	assert.Equal(t, "", ref)
	assert.NotNil(t, err)
	assert.Equal(t, "error validating user: external service error", err.Detail)

	mockUserConnector.AssertExpectations(t)
}

func TestSendMessage_FailsToPushMessage(t *testing.T) {
	mockDB := new(MockDatabase)
	mockUserConnector := new(MockUserConnector)
	service := NewMessageService(mockDB, mockUserConnector)

	mockUserConnector.On("CheckUserExists", "existing_user_1", "Bearer token").Return(true, nil)
	mockUserConnector.On("CheckUserExists", "existing_user_2", "Bearer token").Return(true, nil)
	mockDB.On("SendMessage", "existing_user_1", "existing_user_2", "Hello, World!").
		Return("", goErrors.New("database error"))

	ref, err := service.SendMessage("existing_user_1", "existing_user_2", "Hello, World!", "Bearer token")

	assert.Equal(t, "", ref)
	assert.NotNil(t, err)
	assert.Equal(t, "error sending message: database error", err.Detail)

	mockDB.AssertExpectations(t)
	mockUserConnector.AssertExpectations(t)
}
