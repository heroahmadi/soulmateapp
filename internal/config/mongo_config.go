package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Ctx context.Context

func InitMongoClient() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	var err error
	Client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	Ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := Client.Ping(Ctx, nil); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
}

func CloseMongoClient() {
	err := Client.Disconnect(Ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("MongoDB connection closed")
}
