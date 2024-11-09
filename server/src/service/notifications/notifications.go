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

func (ns NotificationService) SendNotification(receiverId, title, body, authHeader string) *modelErrors.MessageError {

	receiverExists, err := ns.usersConnector.CheckUserExists(receiverId, authHeader)
	if err != nil {
		return modelErrors.ExternalServiceError("error checking user existence: " + err.Error())
	}

	if !receiverExists {
		return modelErrors.ValidationError("receiver does not exist")
	}

	devicesTokens, err := ns.devicesDB.GetDevicesTokens(receiverId)
	if err != nil {
		return modelErrors.InternalServerError("error getting devices tokens: " + err.Error())
	}

	if err := ns.fbConnector.SendNotificationToUserDevices(devicesTokens, title, body); err != nil {
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
