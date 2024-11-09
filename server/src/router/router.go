package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	usersConnector "messages/src/connectors/users-connector"
	messagesController "messages/src/controller/messages"
	notificationsController "messages/src/controller/notifications"
	"messages/src/middleware"
	devicesRepository "messages/src/repository/devices"
	messagesRepository "messages/src/repository/messages"
	messageService "messages/src/service/messages"
	notificationsService "messages/src/service/notifications"
	"os"
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
	var messagesDB messagesRepository.RealTimeDatabaseInterface
	var usersConn usersConnector.Interface
	var devicesDB devicesRepository.DevicesDatabaseInterface

	switch config {
	case MOCK_EXTERNAL:
		messagesDB = messagesRepository.NewMockRealTimeDatabase()
		usersConn = usersConnector.NewMockConnector()
		log.Println("Mocking external connections")
		devicesDB = devicesRepository.NewMockDevicesDatabase()

	default:
		messagesDB = messagesRepository.NewRealTimeDatabase()
		usersConn = usersConnector.NewUsersConnector()
		postgresDB, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			return nil, fmt.Errorf("error connecting to database: %v", err)

		}
		devicesDB, err = devicesRepository.NewDevicesPersistentDatabase(postgresDB)
		if err != nil {
			return nil, fmt.Errorf("error preparing notifications database: %v", err)
		}
	}

	messageServ := messageService.NewMessageService(messagesDB, devicesDB, usersConn)
	messagesContr := messagesController.NewMessageController(messageServ)

	notificationServ := notificationsService.NewNotificationService(devicesDB, usersConn)
	notificationsContr := notificationsController.NewNotificationsController(usersConn, devicesDB, notificationServ)

	private := r.Group("/")
	private.Use(middleware.AuthMiddleware())
	{
		private.GET("/messages", messagesContr.GetMessages)
		private.POST("/messages", messagesContr.SendMessage)
		private.POST("/device", notificationsContr.PostDevice)
		private.POST("/notification", notificationsContr.SendNotification)
	}

	return r, nil

}
