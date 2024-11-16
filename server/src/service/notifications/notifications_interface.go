package service

import "messages/src/model/errors"

type NotificationsServiceInterface interface {
	SendNewMessageNotification(receiverId, senderId, content, chatReference string) *modelErrors.MessageError
	SendFollowerMilestoneNotification(userId, followerId, authHeader string) *modelErrors.MessageError
}
