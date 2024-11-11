package messages

import (
	"log"
	"messages/src/auth"
	"messages/src/model"
	"messages/src/model/errors"
	service "messages/src/service/messages"
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
		modelErrors.SendErrorMessage(ctx, modelErrors.BadRequestError("Error Binding Request: "+err.Error()))
		return
	}
	log.Println("Received Message Request: ", req)

	bearerToken := ctx.GetHeader("Authorization")
	token := strings.Split(bearerToken, "Bearer ")[1]

	claims, err := auth.ValidateToken(token)
	if err != nil {
		modelErrors.SendErrorMessage(ctx, modelErrors.AuthenticationError("Bad token reached controller"))
		return
	}
	senderId := claims.UserId

	ref, er := mc.MessageService.SendMessage(senderId, req.ReceiverId, req.Content, authHeader)
	if er != nil {
		modelErrors.SendErrorMessage(ctx, er)
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
		modelErrors.SendErrorMessage(ctx, err)
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

func (mc *MessageController) GetChatWithUser(context *gin.Context) {

	userId1 := context.Param("userId")
	userId2 := context.GetString("session_user_id")
	authHeader := context.GetHeader("Authorization")

	chat, err := mc.MessageService.GetChatWithUser(userId1, userId2, authHeader)
	if err != nil {
		modelErrors.SendErrorMessage(context, err)
		return
	}

	if chat == nil {
		context.JSON(http.StatusNoContent, gin.H{})
	}

	context.JSON(http.StatusOK, chat)
}
