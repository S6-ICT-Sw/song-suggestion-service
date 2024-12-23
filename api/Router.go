package api

import(
	"song-suggestion-service/api/song_suggestion"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/suggestions", SongSuggestionHandler.)
}