package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbConnected = false
var client *mongo.Client

func ConnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(EnvMongoURI())
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB Atlas")
	dbConnected = true
}

func GetCollection(collectionName string) *mongo.Collection {
	if !dbConnected {
		ConnectDB()
		dbConnected = true
	}
	return client.Database(EnvDBName()).Collection(collectionName)
}