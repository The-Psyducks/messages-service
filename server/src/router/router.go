package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	firebaseConnector "messages/src/connectors/firebase-connector"
	usersConnector "messages/src/connectors/users-connector"
	controllerMessages "messages/src/controller/messages"
	controllerNotifications "messages/src/controller/notifications"
	"messages/src/middleware"
	repositoryDevices "messages/src/repository/devices"
	repositoryMessages "messages/src/repository/messages"
	serviceMessages "messages/src/service/messages"
	serviceNotifications "messages/src/service/notifications"
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
	var messagesDB repositoryMessages.RealTimeDatabaseInterface
	var usersConn usersConnector.Interface
	var devicesDB repositoryDevices.DevicesDatabaseInterface

	switch config {
	case MOCK_EXTERNAL:
		messagesDB = repositoryMessages.NewMockRealTimeDatabase()
		usersConn = usersConnector.NewMockConnector()
		log.Println("Mocking external connections")
		devicesDB = repositoryDevices.NewMockDevicesDatabase()

	default:
		messagesDB = repositoryMessages.NewRealTimeDatabase()
		usersConn = usersConnector.NewUsersConnector()
		postgresDB, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			return nil, fmt.Errorf("error connecting to database: %v", err)

		}
		devicesDB, err = repositoryDevices.NewDevicesPersistentDatabase(postgresDB)
		if err != nil {
			return nil, fmt.Errorf("error preparing notifications database: %v", err)
		}
	}

	messageService := serviceMessages.NewMessageService(messagesDB, devicesDB, usersConn)
	messagesController := controllerMessages.NewMessageController(messageService)

	fbConnector := firebaseConnector.NewFirebaseConnector()
	notificationService := serviceNotifications.NewNotificationService(devicesDB, usersConn, fbConnector)
	notificationsController := controllerNotifications.NewNotificationsController(usersConn, devicesDB, notificationService)

	private := r.Group("/")
	private.Use(middleware.AuthMiddleware())
	{
		private.GET("/messages", messagesController.GetMessages)
		private.POST("/messages", messagesController.SendMessage)
		private.POST("/device", notificationsController.PostDevice)
		private.POST("/notification", notificationsController.SendNotification)
	}

	return r, nil

}
