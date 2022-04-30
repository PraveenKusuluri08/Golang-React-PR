package routes

import (
	"github.com/PraveenKusuluri08/golang-jwt-auth/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("user/signup", controllers.SignUp())
	incomingRoutes.POST("user/signin", controllers.SignIn())
}
