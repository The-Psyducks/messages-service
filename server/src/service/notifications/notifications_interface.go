package service

import "messages/src/model/errors"

type NotificationsServiceInterface interface {
	SendNotification(receiverId, title, body, authHeader string) *modelErrors.MessageError
	SendNewMessageNotification(receiverId, senderId, content, chatReference string) *modelErrors.MessageError
	SendFollowerMilestoneNotification(userId, followers, authHeader string) *modelErrors.MessageError
}
