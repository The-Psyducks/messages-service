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

	claims, _ := auth.ValidateToken(token)
	senderId := claims.UserId

	ref, err := mc.MessageService.SendMessage(senderId, req.ReceiverId, req.Content, authHeader)
	if err != nil {
		errors.SendErrorMessage(ctx, err)
		return
	}
	sendMessageDeliveredResponse(ctx, ref)
}

func sendMessageDeliveredResponse(ctx *gin.Context, ref string) {
	data := model.MessageDeliveredResponse{ChatReference: ref}
	ctx.JSON(http.StatusOK, data)
}

func (mc *MessageController) GetMessages(ctx *gin.Context) {
	//bearerToken := ctx.GetHeader("Authorization")
	//token := strings.Split(bearerToken, "Bearer ")[1]
	//
	//claims, _ := auth.ValidateToken(token)
	userId := "1234" //claims.UserId

	if err := mc.MessageService.GetMessages(userId); err != nil {
		errors.SendErrorMessage(ctx, err)
		return
	}
}
