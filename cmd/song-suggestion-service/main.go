package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"song-suggestion-service/api"
	songsuggestion "song-suggestion-service/api/song_suggestion"
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

	// Initialize the router layer
	sr := api.NewRouter(h)

	// Set up the router and register routes
	r := mux.NewRouter()
	sr.RegisterRoutes(r)

	// Start the HTTP server
	log.Println("Server is running on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", r))
}
