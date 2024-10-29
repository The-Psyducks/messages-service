package router

import (
	"log"
	"messages/src/controller"
	"messages/src/repository"
	"messages/src/service"
	usersConnector "messages/src/user-connector"

	"github.com/gin-gonic/gin"
)

type ConfigurationType int

const (
	// Using iota to define the enum values
	MOCK_EXTERNAL ConfigurationType = iota
	DEFAULT
)

func NewRouter(config ConfigurationType) (*gin.Engine, error) {

	gin.SetMode(gin.ReleaseMode)
	log.Println("Creating router...")

	r := gin.Default()
	var db repository.RealTimeDatabaseInterface
	var users usersConnector.ConnectorInterface
	switch config {
	case MOCK_EXTERNAL:
		db = repository.NewMockRealTimeDatabase()
		users = usersConnector.NewMockConnector()
		//case DEFAULT:
		//	db = repository.NewRealTimeDatabase()
		//	users = usersConnector.NewUsersConnector()
	default:
		db = repository.NewRealTimeDatabase()
		users = usersConnector.NewUsersConnector()

	}

	ms := service.NewMessageService(db, users)
	mc := controller.NewMessageController(ms)

	private := r.Group("/")
	//private.Use(middleware.AuthMiddleware())
	{
		private.GET("/messages", mc.GetMessages)
		private.POST("/messages", mc.SendMessage)
	}

	return r, nil

}
