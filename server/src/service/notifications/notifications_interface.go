package service

import "messages/src/model/errors"

type NotificationsServiceInterface interface {
	SendNotification(receiverId, title, body, authHeader string) *modelErrors.MessageError
}
