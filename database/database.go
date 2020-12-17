package database

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client

var UsersCollection *mongo.Collection

// Connect ...
func Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cli, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URL")))
	if err != nil {
		return err
	}

	if err := cli.Connect(ctx); err != nil {
		return err
	}

	UsersCollection = cli.Database("klever-id").Collection("users")
	client = cli

	return Ping()
}

// Stop ...
func Stop() {
	if client != nil {
		client.Disconnect(context.Background())
	}
}

// Ping ...
func Ping() error {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	return client.Ping(ctx, readpref.Primary())
}
