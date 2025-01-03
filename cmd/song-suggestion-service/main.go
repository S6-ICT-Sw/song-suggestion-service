package main

import (
	"context"
	"log"
	"net/http"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"song-suggestion-service/api"
	songsuggestion "song-suggestion-service/api/song_suggestion"
	"song-suggestion-service/messaging"
	"song-suggestion-service/repository"
	"song-suggestion-service/services"
)

func main() {
	// Connect to MongoDB (update the URI as necessary)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb+srv://song-snippets-admin:DQv4P9LXNBQ2xsdb@songsnippets.ci2mt.mongodb.net/?retryWrites=true&w=majority&appName=SongSnippets"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Initialize repository, service, and handler
	collection := client.Database("suggestionDB").Collection("song_suggestions")
	repo := repository.NewSongSuggestionRepository(collection)
	svc := services.NewSongSuggestionService(repo)
	h := songsuggestion.NewSongSuggestionHandler(svc)

	// Initialize RabbitMQ
	rmq, err := messaging.InitRabbitMQ("song_events")
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer rmq.Close()

	// Start the consumer in the handler
	go h.StartConsumer(rmq)

	// Initialize the router layer
	sr := api.NewRouter(h)

	// Set up the router and register routes
	r := mux.NewRouter()
	sr.RegisterRoutes(r)

	// Set up HTTP server
	server := &http.Server{
		Addr:    ":8081",
		Handler: r,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	<-stop
	log.Println("Shutting down server...")

	// Give the consumer some time to finish processing messages before shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Gracefully stop the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown failed: %v", err)
	}

	// Ensure the RabbitMQ connection is properly closed
	log.Println("Server stopped, RabbitMQ connection closed.")
}
