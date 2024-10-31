package service

import (
	"fmt"
	"log"
	"messages/src/model/errors"
	"messages/src/repository"
	usersConnector "messages/src/user-connector"
	"strings"
)

type MessageService struct {
	db    repository.RealTimeDatabaseInterface
	users usersConnector.ConnectorInterface
}

func NewMessageService(db repository.RealTimeDatabaseInterface, users usersConnector.ConnectorInterface) *MessageService {
	return &MessageService{db, users}
}

func (ms *MessageService) SendMessage(senderId string, receiverId string, content string, authHeader string) (string, *errors.MessageError) {
	//validar que el remitente exista
	senderExists, err := ms.users.CheckUserExists(senderId, authHeader)
	if err != nil {
		log.Printf("error validating user: %v\n", err)
		return "", errors.ExternalServiceError("error validating user: " + err.Error())
	}

	if !senderExists {
		return "", errors.ValidationError("sender does not exist")
	}

	receiverExists, err := ms.users.CheckUserExists(receiverId, authHeader)
	if err != nil {
		log.Printf("error validating user: %v\n", err)
		return "", errors.ExternalServiceError("error validating user: " + err.Error())
	}
	if !receiverExists {
		return "", errors.ValidationError("receiver does not exist")
	}

	//validar que el destinatario exista
	ref, err := ms.db.SendMessage(senderId, receiverId, content)
	if err != nil {
		return "", errors.InternalServerError("error sending message: " + err.Error())
	}
	return ref, nil
}

func (ms *MessageService) GetMessages(id string) ([]string, *errors.MessageError) {
	conversations, err := ms.db.GetConversations(id)
	if err != nil {
		return nil, errors.InternalServerError("error getting conversations: " + err.Error())
	}
	userConversations := filterConversations(id, conversations)
	fmt.Println(userConversations)

	return userConversations, nil
}

func filterConversations(id string, conversations []string) []string {
	result := []string{}
	for _, conversation := range conversations {
		if strings.Contains(conversation, id) {
			result = append(result, conversation)
		}
	}
	return result
}

type MessageServiceInterface interface {
	SendMessage(senderId string, receiverId string, content string, authHeader string) (string, *errors.MessageError)
	GetMessages(id string) ([]string, *errors.MessageError)
}
