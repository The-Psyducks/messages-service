package service

import (
	"log"
	"messages/src/model/errors"
	"messages/src/repository"
	usersConnector "messages/src/user-connector"
)

type MessageService struct {
	db repository.RealTimeDatabaseInterface
}

func NewMessageService(db repository.RealTimeDatabaseInterface) *MessageService {
	return &MessageService{db}
}

func (ms *MessageService) SendMessage(senderId string, receiverId string, content string, authHeader string) *errors.MessageError {
	//validar que el remitente exista
	senderExists, err := usersConnector.CheckUserExists(senderId, authHeader)
	if err != nil {
		log.Printf("error validating user: %v\n", err)
		return errors.ExternalServiceError("error validating user: " + err.Error())
	}

	if !senderExists {
		return errors.ValidationError("sender does not exist")
	}

	receiverExists, err := usersConnector.CheckUserExists(receiverId, authHeader)
	if err != nil {
		log.Printf("error validating user: %v\n", err)
		return errors.ExternalServiceError("error validating user: " + err.Error())
	}
	if !receiverExists {
		return errors.ValidationError("receiver does not exist")
	}

	//validar que el destinatario exista
	if err := ms.db.SendMessage(senderId, receiverId, content); err != nil {
		return errors.InternalServerError("error sending message" + err.Error())
	}
	return nil
}

type MessageServiceInterface interface {
	SendMessage(senderId string, receiverId string, content string, authHeader string) *errors.MessageError
}
