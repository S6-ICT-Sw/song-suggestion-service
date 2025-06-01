package repository

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// NewMongoClient creates and connects a new MongoDB client.
func NewMongoClient(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Ping the database to ensure the connection is successful.
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to MongoDB!")
	return client, nil
}

// CloseMongoClient gracefully disconnects the MongoDB client.
func CloseMongoClient(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Printf("Error disconnecting MongoDB client: %v", err)
	} else {
		log.Println("MongoDB client disconnected successfully.")
	}
}
