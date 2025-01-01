package models

type EditSongSuggestion struct {
	Song_ID string `json:"song_id" bson:"song_id"`
	Title   string `json:"title" bson:"title"`
	Artist  string `json:"artist" bson:"artist"`
}
