package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"messages/src/connectors/users-connector"
	modelErrors "messages/src/model/errors"
	devicesService "messages/src/service/devices"
	"messages/src/service/notifications"

	"messages/src/model"

	"messages/src/repository/devices"
)

type NotificationsController struct {
	uc                   users_connector.Interface
	ds                   devicesService.DevicesServiceInterface
	NotificationsService service.NotificationsServiceInterface
}

func NewNotificationsController(
	uc users_connector.Interface,
	db repository.DevicesDatabaseInterface,
	ns service.NotificationsServiceInterface,
) *NotificationsController {
	ds := devicesService.NewDeviceService(uc, db)
	return &NotificationsController{uc, ds, ns}
}

func (nc *NotificationsController) PostDevice(ctx *gin.Context) {
	userId := ctx.GetString("session_user_id")
	tokenString := ctx.GetString("tokenString")

	var req model.NewDeviceRequest
	if err := ctx.BindJSON(&req); err != nil {
		modelErrors.SendErrorMessage(ctx, modelErrors.BadRequestError("Error Binding Request: "+err.Error()))

		return
	}
	fmt.Println("Add Device Request:", req)
	if err := nc.ds.AddDevice(userId, req.DeviceId, "Bearer "+tokenString); err != nil {
		modelErrors.SendErrorMessage(ctx, err)
		return
	}

	ctx.Status(200)
}

func (nc *NotificationsController) SendNotification(ctx *gin.Context) {
	var notificationRequest model.NotificationRequest

	if err := ctx.BindJSON(&notificationRequest); err != nil {
		modelErrors.SendErrorMessage(ctx, modelErrors.BadRequestError("Error Binding Request: "+err.Error()))
		return
	}
	log.Println("Received Notification Request: ", notificationRequest)
	tokenString := ctx.GetString("tokenString")
	authHeader := "Bearer " + tokenString
	err := nc.NotificationsService.SendNotification(notificationRequest.ReceiverId, notificationRequest.Title, notificationRequest.Body, authHeader)
	if err != nil {
		modelErrors.SendErrorMessage(ctx, err)
		return
	}
}

func (nc *NotificationsController) SendFollowerMilestoneNotification(ctx *gin.Context) {
	var notificationRequest model.UserFollowerMilestoneNotificationRequest
	if err := ctx.BindJSON(&notificationRequest); err != nil {
		modelErrors.SendErrorMessage(ctx, modelErrors.BadRequestError("Error Binding Request: "+err.Error()))
		return
	}
	tokenString := ctx.GetString("tokenString")
	authHeader := "Bearer " + tokenString
	if err := nc.NotificationsService.SendFollowerMilestoneNotification(notificationRequest.UserId, notificationRequest.Followers, authHeader); err != nil {
		modelErrors.SendErrorMessage(ctx, err)
	}

	ctx.Status(204)
	return
}
