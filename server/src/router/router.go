package router

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"messages/src/controller"
	"messages/src/middleware"
	"messages/src/repository"
	"messages/src/service"
	usersConnector "messages/src/user-connector"
	"os"

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
	var rtDb repository.RealTimeDatabaseInterface
	var users usersConnector.ConnectorInterface
	switch config {
	case MOCK_EXTERNAL:
		rtDb = repository.NewMockRealTimeDatabase()
		users = usersConnector.NewMockConnector()
		log.Println("Mocking external connections")
		//case DEFAULT:
		//	db = repository.NewRealTimeDatabase()
		//	users = usersConnector.NewUsersConnector()
	default:
		rtDb = repository.NewRealTimeDatabase()
		users = usersConnector.NewUsersConnector()

	}
	postgresDB, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)

	}

	notificationsDB, err := repository.NewDevicesPersistentDatabase(postgresDB)

	ms := service.NewMessageService(rtDb, notificationsDB, users)
	mc := controller.NewMessageController(ms)

	if err != nil {
		return nil, fmt.Errorf("error preparing notifications database: %v", err)
	}
	nc := controller.NewNotificationsController(users, notificationsDB)

	private := r.Group("/")
	private.Use(middleware.AuthMiddleware())
	{
		private.GET("/messages", mc.GetMessages)
		private.POST("/messages", mc.SendMessage)
		private.POST("/device", nc.PostDevice)
		private.POST("/notification", mc.SendNotication)
	}

	return r, nil

}
