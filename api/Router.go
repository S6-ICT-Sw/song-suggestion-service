package api

import (
	songsuggestion "song-suggestion-service/api/song_suggestion"

	"github.com/gorilla/mux"
)

type Router struct {
	Handler *songsuggestion.SongSuggestionHandler
}

func NewRouter(handler *songsuggestion.SongSuggestionHandler) *Router {
	return &Router{
		Handler: handler,
	}
}

func (sr *Router) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/suggestions", sr.Handler.CreateSongSuggestion).Methods("POST")
	//r.HandleFunc("/suggestions/{id}", sr.Handler.GetSongSuggestionByID).Methods("GET")
}
