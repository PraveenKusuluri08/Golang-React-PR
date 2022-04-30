package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/PraveenKusuluri08/golang-jwt-auth/helpers"
	"github.com/PraveenKusuluri08/golang-jwt-auth/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

const dbName = "USERS"
const collectionName = "user_list"

func SignUp() gin.HandlerFunc {

	return func(c *gin.Context) {
		var user models.User
		collection := helpers.GetCollection(dbName, collectionName)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer cancel()

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, models.Respose{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"Error": err.Error()}})
			return
		}

		count, err := helpers.CountDocuments(bson.M{"email": user.Email}, dbName, collectionName)
		defer cancel()
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusBadRequest, models.Respose{Status: http.StatusBadRequest, Message: "Error while checking a document", Data: map[string]interface{}{"Error": "some thing went really wrong "}})
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, models.Respose{Status: http.StatusBadRequest, Message: "Please check email address", Data: map[string]interface{}{"Error": "please check email address-Email already exists please try again with different email address"}})
			return
		}
		user.CreatedAt = time.Now()
		user.IsExists = true
		user.Role = "USER"
		user.ID = primitive.NewObjectID()
		user.Uid = user.ID.Hex()
		user.Password = helpers.PasswordHasher(user.Password)

		// token, _ := helpers.GenerateAuthJWTToken(user)
		tokenData, _ := helpers.GenerateToken(user.Email, user.Role, user.Uid)

		user.Token = tokenData.Token
		user.RefreshToken = tokenData.RefreshToken

		if validateErr := validate.Struct(&user); validateErr != nil {
			c.JSON(http.StatusBadRequest, models.Respose{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"Error": validateErr.Error()}})
			return
		}

		creatdUser, err := collection.InsertOne(ctx, user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Respose{Status: http.StatusInternalServerError, Message: "Error", Data: map[string]interface{}{"Error": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, models.Respose{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": creatdUser}})
	}
}

func SignIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		var actualUser models.User
		collection := helpers.GetCollection(dbName, collectionName)
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			log.Fatal(err)
			return
		}
		filter := bson.M{"email": user.Email}
		err := collection.FindOne(ctx, filter).Decode(&actualUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid login!", "Message": "Requsting email address is not exists Please try to create new account"})
			return
		}
		match := helpers.CompareHashAndPassword(user.Password, actualUser.Password)

		if !match {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid login!", "Message": "Email and password not matches"})
			return
		}
		token, _ := helpers.GenerateToken(user.Email, actualUser.Role, actualUser.Uid)

		helpers.UpdateToken(token.Token, token.RefreshToken, actualUser.Uid)

		c.JSON(http.StatusOK, actualUser)
	}
}
