package controller

import (
	"messages/src/model"
	"messages/src/model/errors"
	"messages/src/service"

	"github.com/gin-gonic/gin"
)

type MessageController struct {
	MessageService *service.MessageService
}

func NewMessageController(messageService *service.MessageService) *MessageController {
	return &MessageController{messageService}
}

func (mc *MessageController) SendMessage(ctx *gin.Context) {

	authHeader := ctx.GetHeader("Authorization")

	var req model.MessageRequest
	if err := ctx.BindJSON(&req); err != nil {
		errors.NewErrorMessage(ctx, err)
		return
	}

	if err := mc.MessageService.SendMessage(req.SenderId, req.ReceiverId, req.Content, authHeader); err != nil {
		errors.NewErrorMessage(ctx, err)
	}
}
