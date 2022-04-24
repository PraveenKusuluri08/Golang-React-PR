package routes

import (
	"github.com/gin-gonic/gin"
)

func UserRoutes(routes *gin.Engine) {
	//TODO:Use middleware for auth user data
	routes.Use()

	// routes.GET("/users", controllers.GetUsers)
	// routes.GET("/user/:id", controllers.GetUser)
	//routes.PUT("/user/update/:id",controllers.updateUser)
}
