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
	var fbConnector firebaseConnector.Interface

	if err := repositoryMessages.BuildFirebaseConfig(); err != nil {
		log.Fatalln("Error building firebase config:", err)
	}

	switch config {
	case MOCK_EXTERNAL:
		log.Println("Database and connectors")
		messagesDB = repositoryMessages.NewMockRealTimeDatabase()
		usersConn = usersConnector.NewMockConnector()
		devicesDB = repositoryDevices.NewMockDevicesDatabase()
		fbConnector = firebaseConnector.NewMockFirebaseConnector()

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
		fbConnector = firebaseConnector.NewFirebaseConnector()
	}

	notificationService := serviceNotifications.NewNotificationService(devicesDB, usersConn, fbConnector)
	notificationsController := controllerNotifications.NewNotificationsController(usersConn, devicesDB, notificationService)

	messageService := serviceMessages.NewMessageService(messagesDB, devicesDB, usersConn, notificationService)
	messagesController := controllerMessages.NewMessageController(messageService)

	private := r.Group("/")
	private.Use(middleware.AuthMiddleware())
	{
		private.GET("/messages", messagesController.GetMessages)
		private.POST("/messages", messagesController.SendMessage)
		private.GET("/messages/:userId", messagesController.GetChatWithUser)
		private.POST("/device", notificationsController.PostDevice)
		private.POST("/notification/followers-milestone", notificationsController.SendFollowerMilestoneNotification)
		private.POST("/notification/mention", notificationsController.SendMentionNotification)

	}

	return r, nil

}
