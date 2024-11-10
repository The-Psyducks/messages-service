package service

import (
	"log"
	usersConnector "messages/src/connectors/users-connector"
	"messages/src/model/errors"
	repository "messages/src/repository/devices"
	messagesRepository "messages/src/repository/messages"
	serviceNotifications "messages/src/service/notifications"

	"strings"
)

type MessageService struct {
	rtDb                 messagesRepository.RealTimeDatabaseInterface
	dDb                  repository.DevicesDatabaseInterface
	users                usersConnector.Interface
	notificationsService serviceNotifications.NotificationsServiceInterface
}

func NewMessageService(
	rtDb messagesRepository.RealTimeDatabaseInterface,
	dDb repository.DevicesDatabaseInterface,
	users usersConnector.Interface,
	notificationsService serviceNotifications.NotificationsServiceInterface,
) *MessageService {
	return &MessageService{
		rtDb,
		dDb,
		users,
		notificationsService,
	}
}

func (ms *MessageService) SendMessage(senderId string, receiverId string, content string, authHeader string) (string, *modelErrors.MessageError) {
	senderExists, err := ms.users.CheckUserExists(senderId, authHeader)
	if err != nil {
		log.Printf("error validating user: %v\n", err)
		return "", modelErrors.ExternalServiceError("error validating user: " + err.Error())
	}

	if !senderExists {
		return "", modelErrors.ValidationError("sender does not exist")
	}

	receiverExists, err := ms.users.CheckUserExists(receiverId, authHeader)
	if err != nil {
		log.Printf("error validating user: %v\n", err)
		return "", modelErrors.ExternalServiceError("error validating user: " + err.Error())
	}
	if !receiverExists {
		return "", modelErrors.ValidationError("receiver does not exist")
	}

	ref, err := ms.rtDb.SendMessage(senderId, receiverId, content)

	if err != nil {
		return "", modelErrors.InternalServerError("error sending message: " + err.Error())
	}
	if err := ms.notificationsService.SendNewMessageNotification(receiverId, senderId, content, ref); err != nil {
		return "", err
	}
	return ref, nil
}

func (ms *MessageService) GetMessages(id string) ([]string, *modelErrors.MessageError) {
	conversations, err := ms.rtDb.GetConversations()
	if err != nil {
		return nil, modelErrors.InternalServerError("error getting conversations: " + err.Error())
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
