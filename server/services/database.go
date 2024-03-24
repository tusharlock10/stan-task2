package services

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var messagesCollection *mongo.Collection

// Init MongoDB Client
func InitMongoDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	mongoClient = client
	messagesCollection = mongoClient.Database("chatapp").Collection("messages")

	fmt.Println("Database connection successful")
}

// insert the message into database
func InsertMessage(msg *Message) {
	messagesCollection.InsertOne(context.Background(), *msg)
}
