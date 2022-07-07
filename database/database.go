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

var DB *mongo.Client

/************************
For connecting to the DB.
************************/
func GetConnection() (*mongo.Client, error) {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error Loading .env file")
		}
	}

	uri := os.Getenv("MONGO_URI")

	if uri == "" {
		log.Println("Mongo URI Required")
	}
	var err error
	if DB == nil {
		DB, err = mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	}

	if err != nil {
		panic("Error while connecting with DB." + err.Error())
	}
	err = DB.Ping(context.Background(), readpref.Primary())
	if err != nil {
		panic("Error while connecting with DB." + err.Error())
	}
	return DB, nil
}

/***********************************
For getting a particular collection.
***********************************/
func GetCollection(dbName string, collectionName string) (*mongo.Collection, error) {
	client, err := GetConnection()
	if err != nil {
		panic("Error while getting collection " + collectionName + " : " + err.Error())
	}

	collection := client.Database(dbName).Collection(collectionName)

	return collection, nil
}
