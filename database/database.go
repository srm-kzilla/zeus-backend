package database

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//GetConnection for connecting to DB
func GetConnection() (*mongo.Client, error) {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error Loading .env file")
		}
	}

	uri := os.Getenv("MONGO_URI")

	if uri == "" {
		log.Fatal("Mongo URI Required")
	}

	db, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))

	if err != nil {
		panic("Error while connecting with DB." + err.Error())
	}
	err = db.Ping(context.Background(), readpref.Primary())
	if err != nil {
		panic("Error while connecting with DB." + err.Error())
	}
	return db, nil
}

//GetCollection for getting collection
func GetCollection(dbName string, collectionName string) (*mongo.Collection, error) {
	client, err := GetConnection()
	if err != nil {
		panic("Error while getting collection " + collectionName + " : " + err.Error())
	}

	collection := client.Database(dbName).Collection(collectionName)

	return collection, nil
}
