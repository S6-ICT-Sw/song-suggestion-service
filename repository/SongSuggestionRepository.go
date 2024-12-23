package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"song-suggestion-service/models"
)

// ErrNoDocuments is returned when no documents match the query.
var ErrNoDocuments = errors.New("no documents found")

type SongSuggestionRepository struct {
	collection *mongo.Collection
}

/*func NewSongSuggestionRepository(client *mongo.Client, dbName string) *SongSuggestionRepository {
	return &SongSuggestionRepository{
		collection: client.Database(dbName).Collection("song_suggestions"),
	}
}*/

func NewSongSuggestionRepository(collection *mongo.Collection) *SongSuggestionRepository {
	return &SongSuggestionRepository{collection: collection}
}

// Need to fix this
func (r *SongSuggestionRepository) UpdateSong(ctx context.Context, song *models.SongSuggestion) error {
	filter := bson.M{"song_id": song.Song_ID}
	update := bson.M{"$set": bson.M{
		"title":  song.Title,
		"artist": song.Artist,
	}}

	// Use UpdateMany if multiple suggestions could match
	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

// Need to fix this
func (r *SongSuggestionRepository) DeleteSong(ctx context.Context, id string) error {
	filter := bson.M{"song_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

func (r *SongSuggestionRepository) CreateSongSuggestion(ctx context.Context, songSuggestion *models.SongSuggestion) (string, error) {
	result, err := r.collection.InsertOne(ctx, songSuggestion)
	if err != nil {
		return "", err
	}

	// Extract the inserted ID and return it as a string.
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("failed to convert inserted ID to ObjectID")
	}
	return oid.Hex(), nil
}
