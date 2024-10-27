package controller

import (
	"log"
	"messages/src/model"
	"messages/src/model/errors"
	"messages/src/service"
	"net/http"

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
	ref, err := mc.MessageService.SendMessage(req.SenderId, req.ReceiverId, req.Content, authHeader)
	if err != nil {
		errors.SendErrorMessage(ctx, err)
	}
	sendMessageDeliveredResponse(ctx, ref)
}

func sendMessageDeliveredResponse(ctx *gin.Context, ref string) {
	data := model.MessageDeliveredResponse{ref}
	ctx.JSON(http.StatusOK, data)
}
