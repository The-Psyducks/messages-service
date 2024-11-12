package repository

import (
	"errors"
)

type MockRealTimeDatabase struct {
}

func (m *MockRealTimeDatabase) GetChats(string) (*map[string]Message, error) {
	//TODO implement me
	panic("implement me")
}

func NewMockRealTimeDatabase() RealTimeDatabaseInterface {
	return &MockRealTimeDatabase{}
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
