package messages

import (
	"errors"
)

type MockRealTimeDatabase struct {
}

func NewMockRealTimeDatabase() RealTimeDatabaseInterface {
	return &MockRealTimeDatabase{}
}

func (m *MockRealTimeDatabase) SendNotificationToUserDevices(devicesTokens []string, title, body string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockRealTimeDatabase) SendNotification(token string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockRealTimeDatabase) SendMessage(senderId, receiverId, content string) (string, error) {
	if content == "error" {
		return "", errors.New("throwing error in mock")
	}
	return "mockMessageRef", nil
}

func (m *MockRealTimeDatabase) GetConversations() ([]string, error) {
	panic("implement me")
}
