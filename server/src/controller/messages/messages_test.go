package messages

import (
	"bytes"
	"encoding/json"
	"log"
	"messages/src/auth"
	"messages/src/model"
	modelErrors "messages/src/model/errors"
	"messages/src/repository/messages"
	service "messages/src/service/messages"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMessageService struct {
	mock.Mock
}

func (m *MockMessageService) GetChatWithUser(userId1 string, userId2 string, authHeader string) (*model.ChatResponse, *modelErrors.MessageError) {
	//TODO implement me
	panic("implement me")
}

func (m *MockMessageService) SendNotification(receiver, title, body string) *modelErrors.MessageError {
	//TODO implement me
	panic("implement me")
}

func (m *MockMessageService) SendMessage(senderId string, receiverId string, content string, authHeader string) (string, *modelErrors.MessageError) {
	args := m.Called(senderId, receiverId, content, authHeader)
	if err := args.Get(0); err != nil {
		return "", err.(*modelErrors.MessageError)
	}
	return "", nil
}

func (m *MockMessageService) GetMessages(id string, authHeader string) ([]*model.ChatResponse, *modelErrors.MessageError) {
	panic("implement me")
}

func TestSendMessage_Success(t *testing.T) {
	//arrange
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	token, err := auth.GenerateToken("123", "mockUserName", false)
	if err != nil {
		log.Panicln("Error generating token: ", err)
	}

	bearerToken := "Bearer " + token

	mockService := new(MockMessageService)
	controller := NewMessageController(mockService)

	reqBody := model.MessageRequest{ReceiverId: "456", Content: "Hello"}
	mockService.On("SendMessage", "123", "456", "Hello", bearerToken).Return(nil)

	jsonData, _ := json.Marshal(reqBody)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/messages", bytes.NewBuffer(jsonData))
	ctx.Request.Header.Set("Authorization", bearerToken)
	ctx.Request.Header.Set("Content-Type", "application/json")

	//act
	controller.SendMessage(ctx)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestSendMessage_BindJSONError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	token, err := auth.GenerateToken("123", "mockUserName", false)
	if err != nil {
		log.Panicln("Error generating token: ", err)
	}

	bearerToken := "Bearer " + token

	mockService := new(MockMessageService)
	controller := NewMessageController(mockService)

	// Invalid JSON body
	ctx.Request = httptest.NewRequest(http.MethodPost, "/messages", bytes.NewBuffer([]byte("{invalid-json")))
	ctx.Request.Header.Set("Authorization", "Bearer token")
	ctx.Request.Header.Set("Authentication", bearerToken)

	controller.SendMessage(ctx)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Error Binding Request")
}

func TestSendMessage_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	token, err := auth.GenerateToken("123", "mockUserName", false)
	if err != nil {
		log.Panicln("Error generating token: ", err)
	}

	bearerToken := "Bearer " + token

	mockService := new(MockMessageService)
	controller := NewMessageController(mockService)

	reqBody := model.MessageRequest{ReceiverId: "456", Content: "Hello"}
	expectedErr := modelErrors.BadRequestError("Service error") // Simulate an error returned by the service
	mockService.On("SendMessage", "123", "456", "Hello", bearerToken).Return(expectedErr)

	// Prepare the request
	jsonData, _ := json.Marshal(reqBody)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/messages", bytes.NewBuffer(jsonData))
	ctx.Request.Header.Set("Authorization", bearerToken)

	// Call the function
	controller.SendMessage(ctx)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Service error")
	mockService.AssertExpectations(t)
}

type RealTimeDatabaseMock struct {
	mock.Mock
}

func (r *RealTimeDatabaseMock) GetChats(p1 string) (*map[string]repository.Message, error) {
	args := r.Called(p1)
	return args.Get(0).(*map[string]repository.Message), args.Error(1)
}

//func (r *RealTimeDatabaseMock) SendNotificationToUserDevices(devicesTokens []string, title, body string) error {
//	//TODO implement me
//	panic("implement me")
//}

func (r *RealTimeDatabaseMock) SendMessage(senderId string, receiverId string, content string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RealTimeDatabaseMock) GetConversations() ([]string, error) {
	args := r.Called()
	if err := args.Get(1); err != nil {
		return nil, err.(error)
	}
	return args.Get(0).([]string), nil
}

type UsersConnectorMock struct {
	mock.Mock
}

func (u *UsersConnectorMock) GetUserNameAndImage(id string, header string) (string, string, error) {
	args := u.Called(id, header)
	return args.String(0), args.String(1), args.Error(2)
}

func (u *UsersConnectorMock) CheckUserExists(id string, header string) (bool, error) {
	//TODO implement me
	panic("implement me")
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

func TestGetMessages(t *testing.T) {
	//arrange
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	token, err := auth.GenerateToken("1234", "mockUserName", false)
	if err != nil {
		log.Panicln("Error generating token: ", err)
	}
	bearerToken := "Bearer " + token

	realTimeDatabaseMock := new(RealTimeDatabaseMock)
	usersConnectorMock := new(UsersConnectorMock)
	dDbMock := new(MockDevicesDatabase)
	messageService := service.NewMessageService(realTimeDatabaseMock, dDbMock, usersConnectorMock, nil)
	mc := NewMessageController(messageService)

	messages := map[string]repository.Message{
		"1234-4321": {Content: "Hola don pepito", From: "1234", To: "4321", Timestamp: "1"},
		"4321-1234": {Content: "Hola don jose", From: "4321", To: "1234", Timestamp: "2"},
	}

	messages2 := map[string]repository.Message{
		"1234-1111": {Content: "a", From: "1234", To: "1111", Timestamp: "1"},
		"1111-1234": {Content: "b", From: "1111", To: "1234", Timestamp: "2"},
	}

	realTimeDatabaseMock.On("GetChats", "1234-4321").Return(&messages, nil)
	realTimeDatabaseMock.On("GetChats", "1111-1234").Return(&messages2, nil)
	realTimeDatabaseMock.On("GetConversations").Return([]string{"1234-4321", "1111-1234", "1111-1231"}, nil)

	usersConnectorMock.On("GetUserNameAndImage", "1234", bearerToken).
		Return("mockUserName1234", "mockUserImage1234", nil)
	usersConnectorMock.On("GetUserNameAndImage", "4321", bearerToken).
		Return("mockUserName4321", "mockUserImage4321", nil)
	usersConnectorMock.On("GetUserNameAndImage", "1111", bearerToken).
		Return("mockUserName1111", "mockUserImage1111", nil)

	ctx.Request = httptest.NewRequest(http.MethodGet, "/messages", nil)
	ctx.Request.Header.Set("Authorization", bearerToken)
	ctx.Request.Header.Set("Content-Type", "application/json")

	//act
	mc.GetMessages(ctx)

	//assert
	var response []model.ChatResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	expectedResult := []model.ChatResponse{
		{
			ChatReference: "1234-4321",
			UserName:      "mockUserName4321",
			UserImage:     "mockUserImage4321",
			LastMessage:   "Hola don jose",
			Date:          "2",
			ToId:          "4321",
		},
		{
			ChatReference: "1111-1234",
			UserName:      "mockUserName1111",
			UserImage:     "mockUserImage1111",
			LastMessage:   "b",
			Date:          "2",
			ToId:          "1111",
		},
	}
	assert.Equal(t, expectedResult, response)

}
