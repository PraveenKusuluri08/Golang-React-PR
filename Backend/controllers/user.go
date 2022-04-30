package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PraveenKusuluri08/golang-jwt-auth/helpers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUsers() {

}

func GetAllUserData(c *gin.Context) {
	collection := helpers.GetCollection(dbName, collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userUID := c.GetString("uid")
	// opts := options.Find().SetProjection(bson.D{{"email", userEmail}})
	cursor, err := collection.Find(ctx, bson.D{{}}, nil)
	defer cancel()
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error occured while getting documents"})
		return
	}
	var users []primitive.M

	for cursor.Next(ctx) {
		var user primitive.M

		//TODO:Get all users except requested user(admin) data
		if err := cursor.Decode(&user); err != nil {
			log.Fatal(err)
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Something went wrong while decoding"})
			return
		}

		fmt.Print(user["uid"] == userUID)
		if userUID != user["uid"] {
			users = append(users, user)
		}

		defer cursor.Close(ctx)

	}
	c.JSON(http.StatusOK, users)
}
