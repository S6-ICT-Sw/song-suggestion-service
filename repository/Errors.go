package repository

import "errors"

var ErrInvalidSongID = errors.New("song_id must be provided to link the suggestion to a song")
