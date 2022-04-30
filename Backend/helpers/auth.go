package helpers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var dbName = "USERS"
var collectionName = "user_list"

func IsEmailExists(email string) (bool, string) {
	collection := GetCollection(dbName, collectionName)
	filter := bson.M{"email": email}
	var user bson.M
	var singleUser []primitive.M
	err := collection.FindOne(context.Background(), filter).Decode(&user)

	singleUser = append(singleUser, user)

	fmt.Println(user)
	if err != nil {
		return false, "No user exists with the requested email address"
	}

	if len(singleUser[len(singleUser)-1]) != 0 {
		return true, ""
	}
	return false, "User Not exists"
}

func CountDocuments(condition bson.M, dbName string, collectionName string) (int64, error) {
	collection := GetCollection(dbName, collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	count, err := collection.CountDocuments(ctx, condition)
	defer cancel()
	if err != nil {
		return -1, err
	}
	return count, nil
}

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")

		err := ""
		if role != "ADMIN" {
			err = "Unauthorized to view this content! you are not admin"
		}
		if err != "" {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err})
			c.Abort()
			return
		}
		c.Next()
	}
}
