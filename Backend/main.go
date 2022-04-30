package main

import (
	"io"
	"log"
	"os"

	"github.com/PraveenKusuluri08/golang-jwt-auth/helpers"
	"github.com/PraveenKusuluri08/golang-jwt-auth/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error(), "Failed to load env files")
	}
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	router := gin.New()
	router.Use(gin.Logger())

	f, _ := os.Create("gin.log")

	gin.DefaultWriter = io.MultiWriter(f)

	router.Use(gin.LoggerWithFormatter(helpers.LogParser))

	router.Use(gin.Recovery())

	helpers.DBConnection()

	routes.AuthRoutes(router)

	router.Use(helpers.EndPoint())
	routes.UserRoutes(router)
	router.Run(":" + port)
}
