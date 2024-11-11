package repository

import "messages/src/model"

type RealTimeDatabaseInterface interface {
	SendMessage(senderId string, receiverId string, content string) (string, error)
	GetConversations() ([]string, error)
	GetChats() (map[string]model.ChatResponse, error)
}
