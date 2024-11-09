package service

import (
	"log"
	"messages/src/model/errors"
	"messages/src/repository/connectors/user-connector"
	"messages/src/repository/devices"
)

type DevicesServiceInterface interface {
	AddDevice(userId, deviceToken, authToken string) *errors.MessageError
}

type DeviceService struct {
	uc usersConnector.usersConnector
	db devices.DevicesDatabaseInterface
}

func NewDeviceService(uc usersConnector.Interface, db devices.DevicesDatabaseInterface) *DeviceService {
	return &DeviceService{uc: uc, db: db}
}

func (d *DeviceService) AddDevice(userId, deviceToken, authHeader string) *errors.MessageError {
	userIsValid, err := d.uc.CheckUserExists(userId, authHeader)
	if err != nil {
		log.Println("Error validating user: ", err)
		return errors.ExternalServiceError("Error validating user: " + err.Error())
	}

	if !userIsValid {
		log.Println("User does not exist")
		return errors.BadRequestError("User does not exist")
	}

	if err = d.db.AddDevice(userId, deviceToken); err != nil {
		return errors.InternalServerError("Error adding device to db: " + err.Error())
	}

	return nil
}
