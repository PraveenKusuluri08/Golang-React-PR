package helpers

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/PraveenKusuluri08/golang-jwt-auth/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var secretAccessKey = os.Getenv("SECRET_KEY")
var refreshTokenSecretKey = os.Getenv("REFRESH_TOKEN_SECRET")

type SignDetails struct {
	Email     string
	Uid       string
	Role      string
	FirstName string
	jwt.StandardClaims
}

func GenerateAuthJWTToken(user models.User) (*models.TokenDetails, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error(), "Failed to load env files")
	}
	var td = &models.TokenDetails{}
	var err error

	var mySigninKey = []byte(secretAccessKey)

	td.AtExpires = time.Now().Add(time.Hour * 1).Unix()
	td.RtExpires = time.Now().Add(time.Hour * 24).Unix()

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["email"] = user.Email
	atClaims["uid"] = user.Uid
	atClaims["role"] = user.Role
	atClaims["expires"] = td.AtExpires

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString(mySigninKey)
	if err != nil {
		return nil, err
	}
	rtClaims := jwt.MapClaims{}
	rtClaims["email"] = user.Email
	rtClaims["uid"] = user.Uid
	rtClaims["expires"] = td.RtExpires

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(refreshTokenSecretKey))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func GenerateToken(email string, firstName string, uid string, role string) (*models.AuthUsedToken, error) {
	var err error
	tokenData := &models.AuthUsedToken{}
	claims := &SignDetails{
		Email:     email,
		Uid:       uid,
		FirstName: firstName,
		Role:      role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		},
	}
	refreshClaims := &SignDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	tokenData.Token, err = at.SignedString([]byte(secretAccessKey))
	tokenData.RefreshToken, err = rt.SignedString([]byte(secretAccessKey))

	if err != nil {
		log.Panic(err)
		return nil, err
	}
	return tokenData, nil
}

func ValidateJwtAuthToken(signedToken string) (claims *SignDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretAccessKey), nil
		},
	)
	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignDetails)
	if !ok {
		msg = "Token is in valid \n " + err.Error()
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "Token is expired \n " + err.Error()
		return
	}
	return claims, msg
}

//TODO:Renew token when the token is expires
func UpdateToken(token string, refreshToken string, uid string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{"token", token})
	updateObj = append(updateObj, bson.E{"refreshtoken", refreshToken})

	updateAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updateObj = append(updateObj, bson.E{"updatedat", updateAt})

	upsert := true

	filter := bson.M{"uid": uid}

	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	collection := GetCollection(dbName, collectionName)

	_,err:=collection.UpdateOne(ctx, filter, bson.D{
		{"$set", updateObj},
	}, &opt)
	defer cancel()
	if err != nil {
		log.Panic(err)
		return
	}
}
