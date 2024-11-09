package service

import modelErrors "messages/src/model/errors"

type MessageServiceInterface interface {
	SendMessage(senderId string, receiverId string, content string, authHeader string) (string, *modelErrors.MessageError)
	GetMessages(id string) ([]string, *modelErrors.MessageError)
}
