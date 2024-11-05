package controller

import (
	"log"
	"messages/src/auth"
	"messages/src/model"
	"messages/src/model/errors"
	"messages/src/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type MessageController struct {
	MessageService service.MessageServiceInterface
}

func NewMessageController(messageService service.MessageServiceInterface) *MessageController {
	return &MessageController{messageService}
}

func (mc *MessageController) SendMessage(ctx *gin.Context) {

	authHeader := ctx.GetHeader("Authorization")

	var req model.MessageRequest
	if err := ctx.BindJSON(&req); err != nil {
		errors.SendErrorMessage(ctx, errors.BadRequestError("Error Binding Request: "+err.Error()))
		return
	}
	log.Println("Received Message Request: ", req)

	bearerToken := ctx.GetHeader("Authorization")
	token := strings.Split(bearerToken, "Bearer ")[1]

	claims, err := auth.ValidateToken(token)
	if err != nil {
		errors.SendErrorMessage(ctx, errors.AuthenticationError("Bad token reached controller"))
		return
	}
	senderId := claims.UserId

	ref, er := mc.MessageService.SendMessage(senderId, req.ReceiverId, req.Content, authHeader)
	if er != nil {
		errors.SendErrorMessage(ctx, er)
		return
	}
	sendMessageDeliveredResponse(ctx, ref)
}

func (mc *MessageController) GetMessages(ctx *gin.Context) {
	bearerToken := ctx.GetHeader("Authorization")
	token := strings.Split(bearerToken, "Bearer ")[1]

	claims, _ := auth.ValidateToken(token)
	userId := claims.UserId
	conversationReferences, err := mc.MessageService.GetMessages(userId)
	if err != nil {
		errors.SendErrorMessage(ctx, err)
		return
	}
	sendGetMessagesResponse(ctx, conversationReferences)
}

func sendGetMessagesResponse(ctx *gin.Context, references []string) {
	data := model.GetMessagesResponse{ChatReferences: references}
	ctx.JSON(http.StatusOK, data)
}

func sendMessageDeliveredResponse(ctx *gin.Context, ref string) {
	data := model.MessageDeliveredResponse{ChatReference: ref}
	ctx.JSON(http.StatusOK, data)
}

func (mc *MessageController) SendNotication(ctx *gin.Context) {
	var notificationRequest model.NotificationRequest
	if err := ctx.BindJSON(&notificationRequest); err != nil {
		errors.SendErrorMessage(ctx, errors.BadRequestError("Error Binding Request: "+err.Error()))
		return
	}
	log.Println("Received Notification Request: ", notificationRequest)
	
	err := mc.MessageService.SendNotification(notificationRequest.ReceiverId, notificationRequest.Title, notificationRequest.Body)
	if err != nil {
		errors.SendErrorMessage(ctx, err)
		return
	}
}
