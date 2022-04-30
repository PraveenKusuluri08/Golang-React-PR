package routes

import (
	"github.com/PraveenKusuluri08/golang-jwt-auth/controllers"
	"github.com/PraveenKusuluri08/golang-jwt-auth/helpers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(routes *gin.Engine) {
	//TODO:Use middleware for auth user data
	routes.Use(helpers.IsAdmin())

	routes.GET("/users", controllers.GetAllUserData)
	// routes.GET("/user/:id", controllers.GetUser)
	//routes.PUT("/user/update/:id",controllers.updateUser)
}
