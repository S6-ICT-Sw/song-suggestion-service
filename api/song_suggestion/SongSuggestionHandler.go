package songsuggestion

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"os"
	"os/signal"

	"fmt"
	"syscall"

	"github.com/gorilla/mux"

	"song-suggestion-service/messaging"
	"song-suggestion-service/models"
	"song-suggestion-service/services"
)

// SongSuggestionHandler defines the handler for song suggestions.
type SongSuggestionHandler struct {
	service *services.SongSuggestionService
}

// NewSongSuggestionHandler creates a new instance of SongSuggestionHandler.
func NewSongSuggestionHandler(service *services.SongSuggestionService) *SongSuggestionHandler {
	return &SongSuggestionHandler{service: service}
}

func (h *SongSuggestionHandler) CreateSongSuggestion(w http.ResponseWriter, r *http.Request) {
	var suggestion models.SongSuggestion
	if err := json.NewDecoder(r.Body).Decode(&suggestion); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call the service layer to create the suggestion
	id, err := h.service.CreateSuggestion(context.Background(), &suggestion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the created suggestion ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (h *SongSuggestionHandler) DeleteSongSuggestion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Call the service layer to delete the suggestion by ID
	err := h.service.DeleteSuggestionByID(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Song suggestion deleted successfully"})
}

func (h *SongSuggestionHandler) UpdateSongSuggestionsBySongID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	songID := vars["id"]

	var editSongSuggestion models.EditSongSuggestion
	if err := json.NewDecoder(r.Body).Decode(&editSongSuggestion); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.UpdateSongSuggestionsBySongID(r.Context(), songID, &editSongSuggestion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(editSongSuggestion)
	w.WriteHeader(http.StatusOK)
}

func (h *SongSuggestionHandler) GetTopArtistsByName(w http.ResponseWriter, r *http.Request) {
	artistName := r.URL.Query().Get("name")
	if artistName == "" {
		http.Error(w, "artist name is required", http.StatusBadRequest)
		return
	}

	results, err := h.service.GetTopArtistsByName(r.Context(), artistName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// StartConsumer starts the consumer for RabbitMQ to process messages asynchronously.
func (h *SongSuggestionHandler) StartConsumer(rmq *messaging.RabbitMQ) {
	// Create a channel to handle graceful shutdown
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	// Start consuming messages in a goroutine
	go func() {
		// Start consuming the messages from RabbitMQ
		msgs, err := rmq.Consume()
		if err != nil {
			log.Fatalf("Failed to start consuming messages: %v", err)
		}

		// Process messages
		for {
			select {
			case msg := <-msgs:
				// Process the message (e.g., create a song suggestion)
				log.Printf("Received message: %s", msg.Body)
				if err := h.processMessage(msg.Body); err != nil {
					log.Printf("Error processing message: %v", err)
				}
			case <-shutdownCh:
				// Gracefully stop the consumer and close the RabbitMQ connection
				log.Println("Shutdown signal received, stopping the consumer.")
				rmq.Close()
				return
			}
		}
	}()
}

// Process the RabbitMQ message (e.g., create a song suggestion from the message).
func (h *SongSuggestionHandler) processMessage(msgBody []byte) error {
	// Parse the incoming message
	var event struct {
		Song_ID string `json:"song_id"`
		Event   string `json:"event"`
		Title   string `json:"title,omitempty"`  // Optional, only for creation
		Artist  string `json:"artist,omitempty"` // Optional, only for creation
	}
	if err := json.Unmarshal(msgBody, &event); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return fmt.Errorf("invalid message format: %v", err)
	}

	// Log the unmarshalled data for debugging
	log.Printf("Unmarshalled event: %+v", event)

	// Validate SongID
	if event.Song_ID == "" {
		return fmt.Errorf("song_id must be provided")
	}

	// Handle the message based on the event type
	switch event.Event {
	case "created":
		// Handle song creation
		suggestion := models.SongSuggestion{
			Song_ID: event.Song_ID,
			Title:   event.Title,
			Artist:  event.Artist,
		}
		_, err := h.service.CreateSuggestion(context.Background(), &suggestion)
		if err != nil {
			log.Printf("Failed to create song suggestion: %v", err)
			return err
		}
		log.Printf("Song suggestion created successfully for %s by %s", event.Title, event.Artist)

	case "deleted":
		// Handle song deletion
		err := h.service.DeleteSuggestionByID(context.Background(), event.Song_ID)
		if err != nil {
			log.Printf("Failed to delete song suggestions: %v", err)
			return err
		}
		log.Printf("Song suggestions deleted successfully for Song ID %s", event.Song_ID)

	default:
		log.Printf("Unknown event type: %s", event.Event)
		return fmt.Errorf("unknown event type: %s", event.Event)
	}

	return nil
}
