package db

import (
	"context"
	"fmt"
	"log"
	"time"
	"todolist/backend/config"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client
var Database *mongo.Database

func Connect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("Connecting to MongoDB...: ", config.App.MongoURI)
	client, err := mongo.Connect(options.Client().ApplyURI(config.App.MongoURI))
	if err != nil {
		log.Fatalf("MongoDB connect error: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB ping error: %v", err)
	}

	Client = client
	Database = client.Database(config.App.MongoDB)
	log.Printf("Connected to MongoDB: %s / %s", config.App.MongoURI, config.App.MongoDB)
}

func Disconnect() {
	if Client == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := Client.Disconnect(ctx); err != nil {
		log.Printf("MongoDB disconnect error: %v", err)
	}
}

func Col(name string) *mongo.Collection {
	return Database.Collection(name)
}
