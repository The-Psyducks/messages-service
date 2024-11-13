package service

import (
	"log"
	usersConnector "messages/src/connectors/users-connector"
	"messages/src/model"
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

func (ms *MessageService) GetChatWithUser(userId1 string, userId2 string, authHeader string) (*model.ChatResponse, *modelErrors.MessageError) {
	userExists, err := ms.users.CheckUserExists(userId2, authHeader)
	if err != nil {
		log.Printf("error validating user: %v\n", err)
		return nil, modelErrors.ExternalServiceError("error validating user: " + err.Error())
	}

	if !userExists {
		return nil, modelErrors.ValidationError("user does not exists")
	}

	userMessages, modelErr := ms.GetMessages(userId1)
	if modelErr != nil {
		return nil, modelErr
	}
	conversation := filterConversations(userId2, userMessages)
	if len(conversation) == 0 {
		return nil, nil
	}

	if len(conversation) > 1 {
		return nil, modelErrors.InternalServerError("error getting chat: more than one conversation found")
	}

	chats, err := ms.rtDb.GetChats(conversation[0])
	if err != nil {
		return nil, modelErrors.InternalServerError("error getting chats: " + err.Error())
	}

	if len(*chats) == 0 {
		return nil, nil
	}

	return ms.newChat(chats, authHeader)

}

func (ms *MessageService) newChat(chats *map[string]messagesRepository.Message, authHeader string) (*model.ChatResponse, *modelErrors.MessageError) {
	latestChatId := ""
	for chatId, chat := range *chats {
		if latestChatId == "" || (*chats)[latestChatId].Timestamp < chat.Timestamp {
			latestChatId = chatId
		}
	}
	userId := (*chats)[latestChatId].From
	userName, userImage, err := ms.getUserNameAndImage(userId, authHeader)
	if err != nil {
		return nil, modelErrors.ExternalServiceError("error getting user image and name: " + err.Error())
	}

	return &model.ChatResponse{
		ChatReference: latestChatId,
		UserName:      userName,
		UserImage:     userImage,
		LastMessage:   (*chats)[latestChatId].Content,
		Date:          (*chats)[latestChatId].Timestamp,
		ToId:          (*chats)[latestChatId].To,
	}, nil

}

func (ms *MessageService) getUserNameAndImage(id string, authHeader string) (string, string, error) {
	return ms.users.GetUserNameAndImage(id, authHeader)

}

func NewMessageService(
	rtDb messagesRepository.RealTimeDatabaseInterface,
	dDb repository.DevicesDatabaseInterface,
	users usersConnector.Interface,
	notificationsService serviceNotifications.NotificationsServiceInterface,
) MessageServiceInterface {
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
