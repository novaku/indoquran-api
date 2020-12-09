package config

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetMongoDB : get config of MongoDB
func GetMongoDB() (*mongo.Database, error) {
	mongoURI := Config.Database.URI
	mongoDB := Config.Database.DBName

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	return client.Database(mongoDB), nil
}
