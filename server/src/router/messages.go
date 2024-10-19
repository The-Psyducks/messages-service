package router

import (
	"log"

	"messages/src/controller"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	log.Println("Creating router...")

	r := gin.Default()

	mc := controller.NewMessageController()

	r.POST("/messages", mc.SendMessage)

	return r

}
