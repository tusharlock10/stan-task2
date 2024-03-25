package services

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var messagesCollection *mongo.Collection

// Init MongoDB Client
func InitMongoDB() error {
	clientOptions := options.Client().ApplyURI("mongodb://mongodb:27017").SetConnectTimeout(10 * time.Second)

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

// get last 100 messages from database
func GetMessages() (*[]Message, error) {
	var messages []Message

	// setting timout of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find().SetSort(bson.D{{Key: "CreatedAt", Value: -1}}).SetLimit(100)
	cursor, err := messagesCollection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}
	if messages == nil {
		messages = []Message{}
	}

	return &messages, nil
}
