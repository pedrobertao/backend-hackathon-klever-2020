package database

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var Database *mongo.Client

func connect() {
	// client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

}
