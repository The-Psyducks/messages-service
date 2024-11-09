package router

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"messages/src/controller"
	"messages/src/middleware"
	usersConnector2 "messages/src/repository/connectors/user-connector"
	"messages/src/repository/devices"
	"messages/src/repository/messages"
	"messages/src/service"
	"os"

	"github.com/gin-gonic/gin"
)

type ConfigurationType int

const (
	MOCK_EXTERNAL ConfigurationType = iota
	DEFAULT
)

func NewRouter(config ConfigurationType) (*gin.Engine, error) {

	gin.SetMode(gin.ReleaseMode)
	log.Println("Creating router...")

	r := gin.Default()
	var messagesDB messages.RealTimeDatabaseInterface
	var usersConnector usersConnector2.Interface
	var devicesDB devices.DevicesDatabaseInterface

	switch config {
	case MOCK_EXTERNAL:
		messagesDB = messages.NewMockRealTimeDatabase()
		usersConnector = usersConnector2.NewMockConnector()
		log.Println("Mocking external connections")
		devicesDB = devices.NewMockDevicesDatabase()

	default:
		messagesDB = messages.NewRealTimeDatabase()
		usersConnector = usersConnector2.NewUsersConnector()
		postgresDB, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			return nil, fmt.Errorf("error connecting to database: %v", err)

		}
		devicesDB, err = devices.NewDevicesPersistentDatabase(postgresDB)
		if err != nil {
			return nil, fmt.Errorf("error preparing notifications database: %v", err)
		}
	}

	messageService := service.NewMessageService(messagesDB, devicesDB, usersConnector)
	messageController := controller.NewMessageController(messageService)

	notificationsController := controller.NewNotificationsController(usersConnector, devicesDB)

	private := r.Group("/")
	private.Use(middleware.AuthMiddleware())
	{
		private.GET("/messages", messageController.GetMessages)
		private.POST("/messages", messageController.SendMessage)
		private.POST("/device", notificationsController.PostDevice)
		private.POST("/notification", messageController.SendNotication)
	}

	return r, nil

}
