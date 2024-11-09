package service

import (
	"log"
	"messages/src/model/errors"
	"messages/src/repository"
	usersConnector "messages/src/user-connector"
	"strings"
)

type MessageService struct {
	rtDb    repository.RealTimeDatabaseInterface
	dDb	 	repository.DevicesDatabaseInterface
	users 	usersConnector.ConnectorInterface
}

func NewMessageService(rtDb repository.RealTimeDatabaseInterface, dDb repository.DevicesDatabaseInterface,users usersConnector.ConnectorInterface) *MessageService {
	return &MessageService{rtDb, dDb, users}
}

func (ms *MessageService) SendMessage(senderId string, receiverId string, content string, authHeader string) (string, *errors.MessageError) {
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

	ref, err := ms.rtDb.SendMessage(senderId, receiverId, content)

	if err != nil {
		return "", errors.InternalServerError("error sending message: " + err.Error())
	}
	return ref, nil
}

func (ms *MessageService) GetMessages(id string) ([]string, *errors.MessageError) {
	conversations, err := ms.rtDb.GetConversations(id)
	if err != nil {
		return nil, errors.InternalServerError("error getting conversations: " + err.Error())
	}
	userConversations := filterConversations(id, conversations)

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

func (ms *MessageService) SendNotification(receiverId, title, body string) *errors.MessageError {
	devicesTokens, err := ms.dDb.GetDevicesTokens(receiverId)
	if err != nil {
		return errors.InternalServerError("error getting devices tokens: " + err.Error())
	}

	if err := ms.rtDb.SendNotificationToUserDevices(devicesTokens, title, body); err != nil {
		return errors.InternalServerError("error sending notification: " + err.Error())
	}
	return nil
}

type MessageServiceInterface interface {
	SendMessage(senderId string, receiverId string, content string, authHeader string) (string, *errors.MessageError)
	GetMessages(id string) ([]string, *errors.MessageError)
	SendNotification(receiver, title, body string) *errors.MessageError
}
