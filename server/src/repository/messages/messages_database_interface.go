package repository

type RealTimeDatabaseInterface interface {
	SendMessage(senderId string, receiverId string, content string) (string, error)
	GetConversations() ([]string, error)
	GetChats(conversationId string) (*map[string]Message, error)
}
