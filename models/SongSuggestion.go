package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type SongSuggestion struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"` // MongoDB ObjectID for the document
	Song_ID string             `json:"song_id" bson:"song_id"`
	Title   string             `json:"title" bson:"title"`
	Artist  string             `json:"artist" bson:"artist"`
}
