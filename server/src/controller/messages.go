package controller

import (
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
	mc.MessageService.SendMessage("<a pimp named slickback id>", "<fito id>", "https://www.youtube.com/watch?v=RZ1eILLVw74")
}
