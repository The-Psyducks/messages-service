package service

import (
	"log"
	usersConnector "messages/src/connectors/users-connector"
	"messages/src/model/errors"
	"messages/src/repository/devices"
)

type DevicesServiceInterface interface {
	AddDevice(userId, deviceToken, authToken string) *modelErrors.MessageError
}

type DeviceService struct {
	uc usersConnector.Interface
	db repository.DevicesDatabaseInterface
}

func NewDeviceService(uc usersConnector.Interface, db repository.DevicesDatabaseInterface) *DeviceService {
	return &DeviceService{uc: uc, db: db}
}

func (d *DeviceService) AddDevice(userId, deviceToken, authHeader string) *modelErrors.MessageError {
	userIsValid, err := d.uc.CheckUserExists(userId, authHeader)
	if err != nil {
		log.Println("Error validating user: ", err)
		return modelErrors.ExternalServiceError("Error validating user: " + err.Error())
	}

	if !userIsValid {
		log.Println("User does not exist")
		return modelErrors.BadRequestError("User does not exist")
	}

	if err = d.db.AddDevice(userId, deviceToken); err != nil {
		return modelErrors.InternalServerError("Error adding device to db: " + err.Error())
	}

	return nil
}
