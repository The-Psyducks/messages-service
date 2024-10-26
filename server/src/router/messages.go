package router

import (
	"log"
	"messages/src/controller"
	"messages/src/service"

	"github.com/gin-gonic/gin"
)

func NewRouter() (*gin.Engine, error) {

	gin.SetMode(gin.ReleaseMode)
	log.Println("Creating router...")

	r := gin.Default()
	ms := service.NewMessageService()
	mc := controller.NewMessageController(ms)

	r.POST("/messages", mc.SendMessage)

	return r, nil

}
