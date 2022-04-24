package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"id" json:"id"`
	Email        string             `json:"email" validate:"required"`
	Password     string             `json:"password" validate:"required"`
	FirstName    string             `json:"firstName" validate:"required"`
	Age          uint8              `json:"age" validate:"gte=10,lte=100"`
	Token        string             `json:"token"`
	RefreshToken string             `json:"refresh_token "`
	IsExists     bool               `json:"isexists"`
	CreatedAt    time.Time          `json:"createdAt"`
	Role         string             `json:"role" validate:"required,eq=ADMIN|eq=USER"`
	Uid          string             `json:"uid" validate:"required"`
	UpdatedAt    time.Time          `json:"updatedAt"`
}
