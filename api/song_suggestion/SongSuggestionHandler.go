package songsuggestion

import (
	"context"
	"encoding/json"
	"net/http"

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
