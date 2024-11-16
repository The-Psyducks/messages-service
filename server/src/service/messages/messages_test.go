package service

import (
	goErrors "errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"messages/src/model"
	modelErrors "messages/src/model/errors"
	repositoryMessages "messages/src/repository/messages"
	"testing"
)

// Mock for RealTimeDatabaseInterface
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) GetChats(p1 string) (*map[string]repositoryMessages.Message, error) {
	args := m.Called(p1)
	return args.Get(0).(*map[string]repositoryMessages.Message), args.Error(1)
}

func (m *MockDatabase) SendMessage(senderId, receiverId, content string) (string, error) {
	args := m.Called(senderId, receiverId, content)
	return args.String(0), args.Error(1)
}

func (m *MockDatabase) GetConversations() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)

}

// Mock for ConnectorInterface
type MockUserConnector struct {
	mock.Mock
}

func (m *MockUserConnector) GetUserNameAndImage(id string, header string) (string, string, error) {
	args := m.Called(id, header)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockUserConnector) CheckUserExists(userID, authHeader string) (bool, error) {
	args := m.Called(userID, authHeader)
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

type MockNotificationsService struct {
	mock.Mock
}

func (m *MockNotificationsService) SendNewMessageNotification(receiverId, senderId, content, chatReference string) *modelErrors.MessageError {
	_ = m.Called(receiverId, senderId, content, chatReference)
	return nil
}

func (m *MockNotificationsService) SendFollowerMilestoneNotification(userId, followers, authHeader string) *modelErrors.MessageError {
	//TODO implement me
	panic("implement me")
}

func TestSendMessage_HappyPath(t *testing.T) {
	mockDB := new(MockDatabase)
	dDbMock := new(MockDevicesDatabase)
	mockUserConnector := new(MockUserConnector)
	expectedRef := "messageRef123"
	mockNotificationService := new(MockNotificationsService)
	service := NewMessageService(mockDB, dDbMock, mockUserConnector, mockNotificationService)

	mockUserConnector.On("CheckUserExists", "existing_user_1", "Bearer token").
		Return(true, nil)
	mockUserConnector.On("CheckUserExists", "existing_user_2", "Bearer token").
		Return(true, nil)
	mockDB.On("SendMessage", "existing_user_1", "existing_user_2", "Hello, World!").
		Return(expectedRef, nil)
	mockNotificationService.On("SendNewMessageNotification", "existing_user_2", "existing_user_1", "Hello, World!", expectedRef).
		Return(nil)

	ref, err := service.SendMessage("existing_user_1", "existing_user_2", "Hello, World!", "Bearer token")

	assert.Nil(t, err)
	assert.Equal(t, expectedRef, ref)

	mockDB.AssertExpectations(t)
	mockUserConnector.AssertExpectations(t)
}

func TestSendMessage_SenderValidationFails(t *testing.T) {
	mockDB := new(MockDatabase)
	dDbMock := new(MockDevicesDatabase)
	mockUserConnector := new(MockUserConnector)
	service := NewMessageService(mockDB, dDbMock, mockUserConnector, nil)

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
	dDbMock := new(MockDevicesDatabase)
	service := NewMessageService(mockDB, dDbMock, mockUserConnector, nil)

	mockUserConnector.On("CheckUserExists", "existing_user_1", "Bearer token").Return(false, modelErrors.ExternalServiceError("external service error"))

	ref, err := service.SendMessage("existing_user_1", "existing_user_2", "Hello, World!", "Bearer token")

	assert.Equal(t, "", ref)
	assert.NotNil(t, err)
	assert.Equal(t, "error validating user: external service error", err.Detail)

	mockUserConnector.AssertExpectations(t)
}

func TestSendMessage_ReceiverValidationFails(t *testing.T) {
	mockDB := new(MockDatabase)
	mockUserConnector := new(MockUserConnector)
	dDbMock := new(MockDevicesDatabase)
	service := NewMessageService(mockDB, dDbMock, mockUserConnector, nil)

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
	dDbMock := new(MockDevicesDatabase)
	service := NewMessageService(mockDB, dDbMock, mockUserConnector, nil)

	mockUserConnector.On("CheckUserExists", "existing_user_1", "Bearer token").Return(true, nil)
	mockUserConnector.On("CheckUserExists", "existing_user_2", "Bearer token").Return(false, modelErrors.ExternalServiceError("external service error"))

	ref, err := service.SendMessage("existing_user_1", "existing_user_2", "Hello, World!", "Bearer token")

	assert.Equal(t, "", ref)
	assert.NotNil(t, err)
	assert.Equal(t, "error validating user: external service error", err.Detail)

	mockUserConnector.AssertExpectations(t)
}

func TestSendMessage_FailsToPushMessage(t *testing.T) {
	mockDB := new(MockDatabase)
	mockUserConnector := new(MockUserConnector)
	dDbMock := new(MockDevicesDatabase)
	service := NewMessageService(mockDB, dDbMock, mockUserConnector, nil)

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

func TestGetMessagesHappyPath(t *testing.T) {
	//arrange
	mockDB := new(MockDatabase)
	mockUserConnector := new(MockUserConnector)
	dDbMock := new(MockDevicesDatabase)
	service := NewMessageService(mockDB, dDbMock, mockUserConnector, nil)

	mockDB.On("GetConversations").Return([]string{"1234-5678", "1111-1234", "9999-9999"}, nil)

	messages := map[string]repositoryMessages.Message{
		"1234-5678": {Content: "Hola don pepito", From: "1234", To: "5678", Timestamp: "1"},
		"5678-1234": {Content: "Hola don jose", From: "5678", To: "1234", Timestamp: "2"},
	}

	messages2 := map[string]repositoryMessages.Message{
		"1234-1111": {Content: "a", From: "1234", To: "1111", Timestamp: "1"},
		"1111-1234": {Content: "b", From: "1111", To: "1234", Timestamp: "2"},
	}

	_ = map[string]repositoryMessages.Message{
		"9999-9999": {Content: "a", From: "9999", To: "9999", Timestamp: "1"},
	}

	mockDB.On("GetChats", "1234-5678").Return(&messages, nil)
	mockDB.On("GetChats", "1111-1234").Return(&messages2, nil)

	mockUserConnector.On("GetUserNameAndImage", "1234").
		Return("mockUserName1234", "mockUserImage1234", nil)
	mockUserConnector.On("GetUserNameAndImage", "5678", "").
		Return("mockUserName5678", "mockUserImage5678", nil)
	mockUserConnector.On("GetUserNameAndImage", "1111", "").
		Return("mockUserName1111", "mockUserImage1111", nil)

	//act
	resources, err := service.GetMessages("1234", "")
	//assert
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	expectedResult := []*model.ChatResponse{
		{
			ChatReference: "1234-5678",
			UserName:      "mockUserName5678",
			UserImage:     "mockUserImage5678",
			LastMessage:   "Hola don jose",
			Date:          "2",
			ToId:          "1234",
		},
		{
			ChatReference: "1111-1234",
			UserName:      "mockUserName1111",
			UserImage:     "mockUserImage1111",
			LastMessage:   "b",
			Date:          "2",
			ToId:          "1234",
		},
	}
	assert.Equal(t, expectedResult, resources)

}

func TestGetChatHappyPath(t *testing.T) {
	//arrange
	messages := map[string]repositoryMessages.Message{
		"mockNewestRef": {
			To:        "userId2",
			From:      "userId1",
			Id:        "mockMessageId",
			Content:   "mockContent",
			Timestamp: "2",
		},
		"mockOlderRef": {
			To:        "userId1",
			From:      "userId2",
			Id:        "mockMessageId",
			Content:   "mockContent",
			Timestamp: "1",
		},
	}

	mockDB := new(MockDatabase)
	mockUserConnector := new(MockUserConnector)
	dDbMock := new(MockDevicesDatabase)
	mockUserConnector.On("CheckUserExists", "userId2", "authHeader").Return(true, nil)
	service := NewMessageService(mockDB, dDbMock, mockUserConnector, nil)

	mockDB.On("GetConversations").Return([]string{"dm-userId1-userId2", "dm-userId1-userId3"}, nil)
	mockDB.On("GetChats", "dm-userId1-userId2").Return(&messages, nil)

	mockUserConnector.On("GetUserNameAndImage", "userId1", "authHeader").Return("userName", "userImage", nil)

	//act
	resources, err := service.GetChatWithUser("userId1", "userId2", "authHeader")
	//assert
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	assert.Equal(t, &model.ChatResponse{
		ChatReference: "dm-userId1-userId2",
		UserName:      "userName",
		UserImage:     "userImage",
		LastMessage:   "mockContent",
		Date:          "2",
		ToId:          "userId2",
	}, resources)

}

func TestGetChat_UserDoesNotExist(t *testing.T) {
	//arrange
	mockDB := new(MockDatabase)
	mockUserConnector := new(MockUserConnector)
	dDbMock := new(MockDevicesDatabase)
	mockUserConnector.On("CheckUserExists", "userId2", "authHeader").Return(false, nil)
	service := NewMessageService(mockDB, dDbMock, mockUserConnector, nil)

	//act
	resources, err := service.GetChatWithUser("userId1", "userId2", "authHeader")
	//assert
	assert.Nil(t, resources)
	assert.NotNil(t, err)
	assert.Equal(t, "user does not exists", err.Detail)
}

func TestGetChat_ConnectorFails(t *testing.T) {
	//arrange
	mockDB := new(MockDatabase)
	mockUserConnector := new(MockUserConnector)
	dDbMock := new(MockDevicesDatabase)
	mockUserConnector.On("CheckUserExists", "userId2", "authHeader").Return(false, modelErrors.ExternalServiceError("external service error"))
	service := NewMessageService(mockDB, dDbMock, mockUserConnector, nil)

	//act
	resources, err := service.GetChatWithUser("userId1", "userId2", "authHeader")
	//assert
	assert.Nil(t, resources)
	assert.NotNil(t, err)
	assert.Equal(t, "error validating user: external service error", err.Detail)
}

func TestGetChat_GetConversationsFails(t *testing.T) {
	//arrange
	mockDB := new(MockDatabase)
	mockUserConnector := new(MockUserConnector)
	dDbMock := new(MockDevicesDatabase)
	mockUserConnector.On("CheckUserExists", "userId2", "authHeader").Return(true, nil)
	mockDB.On("GetConversations").Return([]string{}, goErrors.New("database error"))
	service := NewMessageService(mockDB, dDbMock, mockUserConnector, nil)

	//act
	resources, err := service.GetChatWithUser("userId1", "userId2", "authHeader")
	//assert
	assert.Nil(t, resources)
	assert.NotNil(t, err)
	assert.Equal(t, "error getting conversations: database error", err.Detail)
}

func TestGetChat_GetChatsFails(t *testing.T) {
	messages := map[string]repositoryMessages.Message{
		"mockNewestRef": {
			To:        "userId2",
			From:      "userId1",
			Id:        "mockMessageId",
			Content:   "mockContent",
			Timestamp: "2",
		},
		"mockOlderRef": {
			To:        "userId1",
			From:      "userId2",
			Id:        "mockMessageId",
			Content:   "mockContent",
			Timestamp: "1",
		},
	}
	//arrange
	mockDB := new(MockDatabase)
	mockUserConnector := new(MockUserConnector)
	dDbMock := new(MockDevicesDatabase)
	mockUserConnector.On("CheckUserExists", "userId2", "authHeader").Return(true, nil)
	mockDB.On("GetConversations").Return([]string{"dm-userId1-userId2"}, nil)
	mockDB.On("GetChats", "dm-userId1-userId2").Return(&messages, goErrors.New("database error"))
	service := NewMessageService(mockDB, dDbMock, mockUserConnector, nil)

	//act
	resources, err := service.GetChatWithUser("userId1", "userId2", "authHeader")
	//assert
	assert.Nil(t, resources)
	assert.NotNil(t, err)
	assert.Equal(t, "error getting chats: database error", err.Detail)
}

func TestGetChat_GetUserNameAndImageFails(t *testing.T) {
	//arrange
	messages := map[string]repositoryMessages.Message{
		"mockNewestRef": {
			To:        "userId2",
			From:      "userId1",
			Id:        "mockMessageId",
			Content:   "mockContent",
			Timestamp: "2",
		},
	}

	mockDB := new(MockDatabase)
	mockUserConnector := new(MockUserConnector)
	dDbMock := new(MockDevicesDatabase)
	mockUserConnector.On("CheckUserExists", "userId2", "authHeader").Return(true, nil)
	service := NewMessageService(mockDB, dDbMock, mockUserConnector, nil)

	mockDB.On("GetConversations").Return([]string{"dm-userId1-userId2"}, nil)
	mockDB.On("GetChats", "dm-userId1-userId2").Return(&messages, nil)
	mockUserConnector.On("GetUserNameAndImage", "userId1", "authHeader").Return("", "", modelErrors.ExternalServiceError("external service error"))

	//act
	resources, err := service.GetChatWithUser("userId1", "userId2", "authHeader")
	//assert
	assert.Nil(t, resources)
	assert.NotNil(t, err)
	assert.Equal(t, "error getting user image and name: external service error", err.Detail)
}

func TestGetChat_NoMessagesReturned(t *testing.T) {
	// arrange
	mockDB := new(MockDatabase)
	mockUserConnector := new(MockUserConnector)
	dDbMock := new(MockDevicesDatabase)
	mockUserConnector.On("CheckUserExists", "userId2", "authHeader").Return(true, nil)
	service := NewMessageService(mockDB, dDbMock, mockUserConnector, nil)

	mockDB.On("GetConversations").Return([]string{"dm-userId1-userId2"}, nil)
	mockDB.On("GetChats", "dm-userId1-userId2").Return(&map[string]repositoryMessages.Message{}, nil)

	// act
	resources, err := service.GetChatWithUser("userId1", "userId2", "authHeader")

	// assert
	assert.Nil(t, resources)
	assert.Nil(t, err)
}

func TestGetChat_InternalErrorMultipleConversations(t *testing.T) {
	// arrange
	mockDB := new(MockDatabase)
	mockUserConnector := new(MockUserConnector)
	dDbMock := new(MockDevicesDatabase)
	mockUserConnector.On("CheckUserExists", "userId2", "authHeader").Return(true, nil)
	service := NewMessageService(mockDB, dDbMock, mockUserConnector, nil)

	mockDB.On("GetConversations").Return([]string{"dm-userId1-userId2", "dm-userId1-userId2"}, nil)

	// act
	resources, err := service.GetChatWithUser("userId1", "userId2", "authHeader")

	// assert
	assert.Nil(t, resources)
	assert.NotNil(t, err)
	assert.Equal(t, "error getting chat: more than one conversation found", err.Detail)
}

func TestGetChat_NoConversationsExist(t *testing.T) {
	// arrange
	mockDB := new(MockDatabase)
	mockUserConnector := new(MockUserConnector)
	dDbMock := new(MockDevicesDatabase)
	mockUserConnector.On("CheckUserExists", "userId2", "authHeader").Return(true, nil)
	service := NewMessageService(mockDB, dDbMock, mockUserConnector, nil)

	mockDB.On("GetConversations").Return([]string{}, nil)

	// act
	resources, err := service.GetChatWithUser("userId1", "userId2", "authHeader")

	// assert
	assert.Nil(t, resources)
	assert.Nil(t, err)
}
