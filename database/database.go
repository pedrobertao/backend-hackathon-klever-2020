package database

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

var UsersCollection *mongo.Collection

// Connect ...
func Connect() error {
	cli, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URL")))
	if err != nil {
		return err
	}
	UsersCollection = cli.Database("klever-id").Collection("users")
	client = cli
	return nil
}

// Stop ...
func Stop() {
	if client != nil {
		client.Disconnect(context.Background())
	}
}
