package service

import (
	"messages/src/model"
	modelErrors "messages/src/model/errors"
)

type MessageServiceInterface interface {
	SendMessage(senderId string, receiverId string, content string, authHeader string) (string, *modelErrors.MessageError)
	GetMessages(id string) ([]string, *modelErrors.MessageError)
	GetChatWithUser(userId1 string, userId2 string, authHeader string) (*model.ChatResponse, *modelErrors.MessageError)
}
