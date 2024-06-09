package migration

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"soulmateapp/api/model"
	"soulmateapp/internal/config"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbName string = "soulmate"
var collectionName string = "users"

func InitData() {
	if err := config.Client.Ping(context.Background(), nil); err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
		return
	}
	log.Println("Connected to MongoDB!")

	createUserIndex()
	initUserData()

	log.Println("Migration completed")
}

func createUserIndex() {
	log.Println("Creating index")
	collection := config.Client.Database(dbName).Collection(collectionName)
	indexModel := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	_, err := collection.Indexes().CreateMany(context.Background(), indexModel)
	if err != nil && err.Error() != "index already exists with different options" {
		log.Fatal(err)
	}
}

func initUserData() {
	log.Println("Init user data")
	collection := config.Client.Database(dbName).Collection(collectionName)
	truncateCollection(collection)

	data, err := os.ReadFile("db/user.json")
	if err != nil {
		log.Fatal(err)
	}

	var users []model.User
	if err := json.Unmarshal(data, &users); err != nil {
		log.Fatal(err)
	}

	for _, user := range users {
		id, errUuid := uuid.NewRandom()
		if errUuid != nil {
			fmt.Println("Error generating UUID:", errUuid)
			return
		}
		user.ID = id.String()
		_, err := collection.InsertOne(context.Background(), user)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func truncateCollection(collection *mongo.Collection) {
	_, err := collection.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("collection trucated")
}
