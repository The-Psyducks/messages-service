package router

import (
	"log"
	"messages/src/controller"
	"messages/src/middleware"
	"messages/src/repository"
	"messages/src/service"
	usersConnector "messages/src/user-connector"

	"github.com/gin-gonic/gin"
)

func NewRouter() (*gin.Engine, error) {

	gin.SetMode(gin.ReleaseMode)
	log.Println("Creating router...")

	r := gin.Default()
	db := repository.NewRealTimeDatabase()
	users := usersConnector.NewUsersConnector()
	ms := service.NewMessageService(db, users)
	mc := controller.NewMessageController(ms)

	private := r.Group("/")
	private.Use(middleware.AuthMiddleware())
	{
		private.POST("/messages", mc.SendMessage)
	}

	return r, nil

}

/*private := r.Engine.Group("/")
private.Use(middleware.AuthMiddleware())
{
	private.GET("/users/:id", userController.GetUserProfileById)
	private.PUT("/users/profile", userController.ModifyUserProfile)

	private.POST("/users/:id/follow", userController.FollowUser)
	private.DELETE("/users/:id/follow", userController.UnfollowUser)
	private.GET("/users/:id/followers", userController.GetFollowers)
	private.GET("/users/:id/following", userController.GetFollowing)

	private.GET("/users/search", userController.SearchUsers)

	private.GET("/users/all", userController.GetAllUsers)
}*/
