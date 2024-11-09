package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"messages/src/model"
	"messages/src/model/errors"
	"messages/src/repository/connectors/user-connector"
	"messages/src/repository/devices"
	"messages/src/service"
)

type NotificationsController struct {
	ds service.DevicesServiceInterface
}

func NewNotificationsController(
	uc usersConnector.usersConnector,
	db devices.DevicesDatabaseInterface) *NotificationsController {
	ds := service.NewDeviceService(uc, db)
	return &NotificationsController{ds}
}

func (nc *NotificationsController) PostDevice(ctx *gin.Context) {
	userId := ctx.GetString("session_user_id")
	tokenString := ctx.GetString("tokenString")

	var req model.NewDeviceRequest
	if err := ctx.BindJSON(&req); err != nil {
		errors.SendErrorMessage(ctx, errors.BadRequestError("Error Binding Request: "+err.Error()))

		return
	}
	fmt.Println("Add Device Request:", req)
	if err := nc.ds.AddDevice(userId, req.DeviceId, "Bearer "+tokenString); err != nil {
		errors.SendErrorMessage(ctx, err)
		return
	}

	ctx.Status(200)
}

//func NewNotificationsController(usersConnectorMock, devicesDatabaseMock)
