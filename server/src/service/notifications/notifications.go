package service

import (
	firebaseConnector "messages/src/connectors/firebase-connector"
	usersConnector "messages/src/connectors/users-connector"
	modelErrors "messages/src/model/errors"
	"messages/src/repository/devices"
)

type NotificationService struct {
	devicesDB      repository.DevicesDatabaseInterface
	usersConnector usersConnector.Interface
	fbConnector    firebaseConnector.Interface
}

func (ns *NotificationService) SendMentionNotification(userId string, taggerId string, postId string, authHeader string) *modelErrors.MessageError {
	userExists, err := ns.usersConnector.CheckUserExists(userId, authHeader)
	if err != nil {
		return modelErrors.ExternalServiceError("error checking user existence: " + err.Error())
	}
	if !userExists {
		return modelErrors.ValidationError("receiver does not exist")
	}

	followerExists, err := ns.usersConnector.CheckUserExists(userId, authHeader)
	if err != nil {
		return modelErrors.ExternalServiceError("error checking user existence: " + err.Error())
	}
	if !followerExists {
		return modelErrors.ValidationError("receiver does not exist")
	}

	devicesTokens, err := ns.devicesDB.GetDevicesTokens(userId)
	if err != nil {
		return modelErrors.InternalServerError("error getting devices tokens: " + err.Error())
	}

	data := map[string]string{
		"deeplink": "twitSnap://home_twitSnap?twitSnapId=" + postId,
	}

	if err := ns.fbConnector.SendNotificationToUserDevices(
		devicesTokens,
		"Yay! You got mentioned",
		"Someone mentioned you in a post! (๑˃̵ᴗ˂̵)و",
		data); err != nil {
		return modelErrors.InternalServerError("error sending notification: " + err.Error())
	}
	return nil
}

func NewNotificationService(devicesDB repository.DevicesDatabaseInterface, usersConnector usersConnector.Interface, fbConnector firebaseConnector.Interface) NotificationsServiceInterface {
	return &NotificationService{
		devicesDB:      devicesDB,
		usersConnector: usersConnector,
		fbConnector:    fbConnector,
	}
}

func (ns *NotificationService) SendNewMessageNotification(receiverId string, senderId string, content string, chatReference string) *modelErrors.MessageError {

	devicesTokens, err := ns.devicesDB.GetDevicesTokens(receiverId)
	if err != nil {
		return modelErrors.InternalServerError("error getting devices tokens: " + err.Error())
	}

	data := map[string]string{
		"deeplink": "twitSnap://messages_chat?userId=" + senderId + "?refId=" + chatReference,
		//twitSnap://messages_chat?userId=123?refId=asdvas
	}
	if err := ns.fbConnector.SendNotificationToUserDevices(devicesTokens, "New message", content, data); err != nil {
		return modelErrors.InternalServerError("error sending notification: " + err.Error())
	}
	return nil
}

func (ns *NotificationService) SendFollowerMilestoneNotification(userId, followerId, authHeader string) *modelErrors.MessageError {
	userExists, err := ns.usersConnector.CheckUserExists(userId, authHeader)
	if err != nil {
		return modelErrors.ExternalServiceError("error checking user existence: " + err.Error())
	}
	if !userExists {
		return modelErrors.ValidationError("receiver does not exist")
	}

	followerExists, err := ns.usersConnector.CheckUserExists(userId, authHeader)
	if err != nil {
		return modelErrors.ExternalServiceError("error checking user existence: " + err.Error())
	}
	if !followerExists {
		return modelErrors.ValidationError("receiver does not exist")
	}

	devicesTokens, err := ns.devicesDB.GetDevicesTokens(userId)
	if err != nil {
		return modelErrors.InternalServerError("error getting devices tokens: " + err.Error())
	}

	data := map[string]string{
		"deeplink": "twitSnap://profile_profile?userId=" + followerId,
	}

	if err := ns.fbConnector.SendNotificationToUserDevices(
		devicesTokens,
		"New follower!!",
		"You got yourself a new follower (˶ᵔ ᵕ ᵔ˶) ",
		data); err != nil {
		return modelErrors.InternalServerError("error sending notification: " + err.Error())
	}
	return nil
}
