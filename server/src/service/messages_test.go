package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for RealTimeDatabaseInterface
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) SendMessage(senderId, receiverId, content string) (string, error) {
	args := m.Called(senderId, receiverId, content)
	return args.String(0), args.Error(1)
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
	// Create mocks
	mockDB := new(MockDatabase)
	mockUserConnector := new(MockUserConnector)

	// Expected values
	expectedRef := "messageRef123"

	// Initialize the MessageService with mocks
	service := NewMessageService(mockDB, mockUserConnector)

	// Set up expectations for mocks
	mockUserConnector.On("CheckUserExists", "existing_user_1", "Bearer token").Return(true, nil)
	mockUserConnector.On("CheckUserExists", "existing_user_2", "Bearer token").Return(true, nil)
	mockDB.On("SendMessage", "existing_user_1", "existing_user_2", "Hello, World!").Return(expectedRef, nil)

	// Execute the SendMessage function
	ref, err := service.SendMessage("existing_user_1", "existing_user_2", "Hello, World!", "Bearer token")

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, expectedRef, ref)

	// Verify expectations
	mockDB.AssertExpectations(t)
	mockUserConnector.AssertExpectations(t)
}
