package songsuggestion

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

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
