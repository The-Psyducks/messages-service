package controller

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"messages/src/auth"
	"messages/src/repository/connectors/user-connector"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockConnector struct {
	mock.Mock
}

func (m *mockConnector) CheckUserExists(userId, deviceToken string) (bool, error) {
	args := m.Called(userId, deviceToken)
	return args.Bool(0), args.Error(1)
}

var _ usersConnector.usersConnector = (*mockConnector)(nil)

type mockDatabase struct {
	mock.Mock
}

func (m *mockDatabase) GetDevicesTokens(id string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockDatabase) AddDevice(id string, token string) error {
	_ = m.Called(id, token)
	return nil
}

func TestAddDeviceForUser(t *testing.T) {
	//arrange
	usersConnectorMock := new(mockConnector)
	devicesDatabaseMock := new(mockDatabase)
	token, _ := auth.GenerateToken("userId", "username", false)
	bearerToken := "Bearer " + token
	usersConnectorMock.On("CheckUserExists", "userId", bearerToken).Return(true, nil)
	devicesDatabaseMock.On("AddDevice", "userId", "deviceToken").Return(nil)

	_ = gin.Default()
	// Create a new gin context
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	body := `{"device_id": "deviceToken"}`
	ctx.Request, _ = http.NewRequest("POST", "/device", bytes.NewBufferString(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Request.Header.Set("Authorization", bearerToken)
	ctx.Set("session_user_id", "userId")
	ctx.Set("tokenString", token)

	nc := NewNotificationsController(usersConnectorMock, devicesDatabaseMock)

	//act
	nc.PostDevice(ctx)

	//assert
	usersConnectorMock.AssertCalled(t, "CheckUserExists", "userId", bearerToken)
	devicesDatabaseMock.AssertCalled(t, "AddDevice", "userId", "deviceToken")
}
