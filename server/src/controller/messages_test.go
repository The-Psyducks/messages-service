package controller

import (
	"bytes"
	"encoding/json"
	"log"
	"messages/src/auth"
	"messages/src/model"
	"messages/src/model/errors"
	"messages/src/service"
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

func (m *MockMessageService) SendNotification(receiver, title, body string) *errors.MessageError {
	//TODO implement me
	panic("implement me")
}

func (m *MockMessageService) SendMessage(senderId string, receiverId string, content string, authHeader string) (string, *errors.MessageError) {
	args := m.Called(senderId, receiverId, content, authHeader)
	if err := args.Get(0); err != nil {
		return "", err.(*errors.MessageError)
	}
	return "", nil
}

func (m *MockMessageService) GetMessages(id string) ([]string, *errors.MessageError) {
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
	expectedErr := errors.BadRequestError("Service error") // Simulate an error returned by the service
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

func (r *RealTimeDatabaseMock) SendNotificationToUserDevices(devicesTokens []string, title, body string) error {
	//TODO implement me
	panic("implement me")
}

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
	messageService := service.NewMessageService(realTimeDatabaseMock, dDbMock, usersConnectorMock)
	mc := NewMessageController(messageService)

	ctx.Request = httptest.NewRequest(http.MethodGet, "/messages", nil)
	ctx.Request.Header.Set("Authorization", bearerToken)
	ctx.Request.Header.Set("Content-Type", "application/json")

	realTimeDatabaseMock.On("GetConversations").Return([]string{"1234-4321", "1111-1234", "1111-1231"}, nil)
	//act
	mc.GetMessages(ctx)

	//assert
	var response model.GetMessagesResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, []string{"1234-4321", "1111-1234"}, response.ChatReferences)

}
