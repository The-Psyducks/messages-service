package controller

import "github.com/gin-gonic/gin"

type MessageController struct{}

func NewMessageController() *MessageController {
	return nil
}

func (mc *MessageController) SendMessage(ctx *gin.Context) {

}
