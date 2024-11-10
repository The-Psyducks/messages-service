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
		"senderId": senderId,
		"deeplink": "twitSnap://messages_chat?refId=" + chatReference,
	}
	if err := ns.fbConnector.SendNotificationToUserDevices(devicesTokens, "New message", content, data); err != nil {
		return modelErrors.InternalServerError("error sending notification: " + err.Error())
	}
	return nil
}

func (ns *NotificationService) SendFollowerMilestoneNotification(userId string, followers string, authHeader string) *modelErrors.MessageError {
	receiverExists, err := ns.usersConnector.CheckUserExists(userId, authHeader)
	if err != nil {
		return modelErrors.ExternalServiceError("error checking user existence: " + err.Error())
	}

	if !receiverExists {
		return modelErrors.ValidationError("receiver does not exist")
	}

	devicesTokens, err := ns.devicesDB.GetDevicesTokens(userId)
	if err != nil {
		return modelErrors.InternalServerError("error getting devices tokens: " + err.Error())
	}

	data := map[string]string{
		"deeplink": "twitSnap://messages_chat?refId=" + userId,
	}

	if err := ns.fbConnector.SendNotificationToUserDevices(
		devicesTokens,
		"New milestone!!",
		"You just reached "+followers+" followers!! (˶ᵔ ᵕ ᵔ˶) ",
		data); err != nil {
		return modelErrors.InternalServerError("error sending notification: " + err.Error())
	}
	return nil
}
