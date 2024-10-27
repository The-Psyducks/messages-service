package repository

import "errors"

type MockRealTimeDatabase struct {
}

func NewMockRealTimeDatabase() RealTimeDatabaseInterface {
	return &MockRealTimeDatabase{}
}

func (m *MockRealTimeDatabase) SendMessage(senderId, receiverId, content string) (string, error) {
	if content == "error" {
		return "", errors.New("throwing error in mock")
	}
	return "mockMessageRef", nil
}
