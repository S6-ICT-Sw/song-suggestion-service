package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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
func (r *SongSuggestionRepository) UpdateBySongID(ctx context.Context, songID string, editSong *models.EditSongSuggestion) error {
	filter := bson.M{"song_id": songID}
	edit := bson.M{"$set": editSong}
	result, err := r.collection.UpdateOne(ctx, filter, edit)
	if err != nil {
		return err // Error formatting is handled in the handler layer.
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *SongSuggestionRepository) DeleteBySongID(ctx context.Context, songID string) error {
	filter := bson.M{"song_id": songID}
	result, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return err // Error formatting is handled in the handler layer.
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
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

func (r *SongSuggestionRepository) GetTopArtistsByName(ctx context.Context, artistName string) ([]models.SongSuggestion, error) {
	filter := bson.M{
		"artist": bson.M{
			"$regex":   artistName, // Partial match
			"$options": "i",        // Case-insensitive
		},
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetLimit(5))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.SongSuggestion
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
