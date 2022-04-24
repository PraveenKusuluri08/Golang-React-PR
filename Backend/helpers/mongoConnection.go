package helpers

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBConnection() *mongo.Client {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error(), "Failed to load env files")
	}
	MongoDB := os.Getenv("MONGO_URI")

	clientOptions := options.Client().ApplyURI(MongoDB)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("INFO: DataBase connected😍", log.Ldate|log.Ltime)
	return client
}

var Client *mongo.Client = DBConnection()

func GetCollection(dbName string, collectionName string) *mongo.Collection {
	collection := Client.Database(dbName).Collection(collectionName)
	return collection
}
