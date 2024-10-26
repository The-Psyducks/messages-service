package service

import (
	"errors"
	"messages/src/repository"
	usersConnector "messages/src/user-connector"
)

type MessageService struct {
}

func NewMessageService() *MessageService {
	return &MessageService{}
}

func (ms *MessageService) SendMessage(senderId string, receiverId string, content string, authHeader string) error {
	//validar que el remitente exista
	senderExists, err := usersConnector.CheckUserExists(senderId, authHeader)
	if err != nil {
		return err
	}

	if !senderExists {
		return errors.New("sender does not exist")
	}

	receiverExists, err := usersConnector.CheckUserExists(receiverId, authHeader)
	if err != nil {
		return err
	}
	if !receiverExists {
		return errors.New("receiver does not exist")
	}

	//validar que el destinatario exista
	if err := repository.SendMessage(senderId, receiverId, content); err != nil {
		return err
	}
	return nil
}
