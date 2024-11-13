package repository

import (
	"errors"
)

type MockRealTimeDatabase struct {
}

func (m *MockRealTimeDatabase) GetChats(behaviour string) (*map[string]Message, error) {
	//switch behaviour {
	//case "ok":
	//	return &map[string]Message{
	//		"mockNewestRef": Message{
	//			To:        "mockReceiverId",
	//			From:      "mockSenderId",
	//			Id:        "mockMessageId",
	//			Content:   "mockContent",
	//			Timestamp: "2",
	//		},
	//		"mockOlderRef": Message{
	//			To:        "mockReceiverId",
	//			From:      "mockSenderId",
	//			Id:        "mockMessageId",
	//			Content:   "mockContent",
	//			Timestamp: "1",
	//		},
	//	}, nil
	//case "error":
	//	return nil, errors.New("throwing error in mock")
	//default:
	//	panic("behaviour should be ok or error but it was: " + behaviour)
	//}

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
