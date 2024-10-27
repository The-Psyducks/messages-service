package controller

import (
	"bytes"
	"encoding/json"
	"messages/src/model"
	"messages/src/model/errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock of MessageServiceInterface
type MockMessageService struct {
	mock.Mock
}

func (m *MockMessageService) SendMessage(senderId string, receiverId string, content string, authHeader string) (string, *errors.MessageError) {
	args := m.Called(senderId, receiverId, content, authHeader)
	if err := args.Get(0); err != nil {
		return "", err.(*errors.MessageError)
	}
	return "", nil
}

func TestSendMessage_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Prepare mock service and controller
	mockService := new(MockMessageService)
	controller := NewMessageController(mockService)

	// Define input and expected behavior
	reqBody := model.MessageRequest{SenderId: "123", ReceiverId: "456", Content: "Hello"}
	mockService.On("SendMessage", "123", "456", "Hello", "Bearer token").Return(nil)

	// Prepare the request
	jsonData, _ := json.Marshal(reqBody)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/send", bytes.NewBuffer(jsonData))
	ctx.Request.Header.Set("Authorization", "Bearer token")

	// Call the function
	controller.SendMessage(ctx)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestSendMessage_BindJSONError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	mockService := new(MockMessageService)
	controller := NewMessageController(mockService)

	// Simulate invalid JSON body
	ctx.Request = httptest.NewRequest(http.MethodPost, "/send", bytes.NewBuffer([]byte("{invalid-json")))
	ctx.Request.Header.Set("Authorization", "Bearer token")

	controller.SendMessage(ctx)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Error Binding Request")
}

func TestSendMessage_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	mockService := new(MockMessageService)
	controller := NewMessageController(mockService)

	// Define input and expected error from the service
	reqBody := model.MessageRequest{SenderId: "123", ReceiverId: "456", Content: "Hello"}
	expectedErr := errors.BadRequestError("Service error") // Simulate an error returned by the service
	mockService.On("SendMessage", "123", "456", "Hello", "Bearer token").Return(expectedErr)

	// Prepare the request
	jsonData, _ := json.Marshal(reqBody)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/send", bytes.NewBuffer(jsonData))
	ctx.Request.Header.Set("Authorization", "Bearer token")

	// Call the function
	controller.SendMessage(ctx)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Service error")
	mockService.AssertExpectations(t)
}
