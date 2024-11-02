package controller

import (
	"github.com/gin-gonic/gin"
	"messages/src/repository"
	usersConnector "messages/src/user-connector"
)

type NotificationsController struct {
	uc usersConnector.ConnectorInterface
	db repository.DevicesDatabaseInterface
}

func NewNotificationsController(
	uc usersConnector.ConnectorInterface,
	db repository.DevicesDatabaseInterface) *NotificationsController {
	return &NotificationsController{uc: uc, db: db}
}

func (c *NotificationsController) PostDevice(ctx *gin.Context) {

}

//func NewNotificationsController(usersConnectorMock, devicesDatabaseMock)
