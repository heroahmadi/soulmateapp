package migration

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"soulmateapp/api/model/entity"
	"soulmateapp/internal/config"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var dbName string = "soulmate"
var usersCollectionName string = "users"
var userLikesCollectionName string = "user_likes"

func InitData() {
	if err := config.Client.Ping(context.Background(), nil); err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
		return
	}
	log.Println("Connected to MongoDB!")

	createUserIndex()
	initUserData()
	createUserLikesIndex()
	initUserLikesData()

	log.Println("Migration completed")
}

func createUserIndex() {
	log.Println("Creating index")
	collection := config.Client.Database(dbName).Collection(usersCollectionName)
	indexModel := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "username", Value: 1}, {Key: "id", Value: 1}},
		},
	}
	_, err := collection.Indexes().CreateMany(context.Background(), indexModel)
	if err != nil && err.Error() != "index already exists with different options" {
		log.Fatal(err)
	}
}

func initUserData() {
	log.Println("Init user data")
	collection := config.Client.Database(dbName).Collection(usersCollectionName)
	truncateCollection(collection)

	data, err := os.ReadFile("db/user.json")
	if err != nil {
		log.Fatal(err)
	}

	var users []entity.User
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

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}
		user.Password = string(hashedPassword)
		_, err = collection.InsertOne(context.Background(), user)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func truncateCollection(collection *mongo.Collection) {
	_, err := collection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("collection trucated")
}

func createUserLikesIndex() {
	log.Println("Creating index user likes")
	collection := config.Client.Database(dbName).Collection(userLikesCollectionName)
	indexModel := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "date", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	_, err := collection.Indexes().CreateMany(context.Background(), indexModel)
	if err != nil && err.Error() != "index already exists with different options" {
		log.Fatal(err)
	}
}

func initUserLikesData() {
	log.Println("Init user likes data")
	collection := config.Client.Database(dbName).Collection(userLikesCollectionName)
	truncateCollection(collection)
}
