package services

import (
	"context"
	"errors"

	"song-suggestion-service/models"
	"song-suggestion-service/repository"
)

type SongSuggestionService struct {
	repo repository.SongSuggestionRepository // Dependency injection of repository
}

func NewSongSuggestionService(repo repository.SongSuggestionRepository) *SongSuggestionService {
	return &SongSuggestionService{repo: repo}
}

func (s *SongSuggestionService) CreateSuggestion(ctx context.Context, suggestion *models.SongSuggestion) (string, error) {
	if suggestion.Song_ID == "" {
		return "", errors.New("song_id must be provided")
	}

	// Call the repository to create the suggestion
	id, err := s.repo.CreateSongSuggestion(ctx, suggestion)
	if err != nil {
		return "", err // Propagate repository error to handler
	}
	return id, nil
}
