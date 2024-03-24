package services

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var messagesCollection *mongo.Collection

// Init MongoDB Client
func InitMongoDB() error {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017").SetConnectTimeout(10 * time.Second)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}
	if err := client.Ping(context.Background(), nil); err != nil {
		return err
	}

	mongoClient = client
	messagesCollection = mongoClient.Database("chatapp").Collection("messages")

	fmt.Println("Database connection successful")
	return nil
}

// insert the message into database
func InsertMessage(msg *Message) {
	messagesCollection.InsertOne(context.Background(), *msg)
}
