package repositories

import (
	"spotigram/internal/customerrors"
	"spotigram/internal/infrastructure/abstractions"

	"github.com/gocql/gocql"
)

type CqlPlaylistSongRepository struct {
	DBProvider abstractions.CqlDatabaseProvider
}

// Adds a message to the repository.
// May return ErrInternal on failure.
func (cps *CqlPlaylistSongRepository) AddPlaylistSong(playlistId, songId string) error {
	session := cps.DBProvider.GetSession()

	insertStmt := session.Query(`
		INSERT INTO playlist_songs (playlist_id, song_id)
		VALUES (?, ?)
	`)

	err := insertStmt.Bind(
		playlistId,
		songId,
	).Exec()

	if err != nil {
		return &customerrors.ErrInternal{Message: err.Error()}
	}

	return nil
}

// Adds a message to the repository.
// May return ErrInternal on failure.
func (cps *CqlPlaylistSongRepository) DeletePlaylistSong(playlistId, songId string) error {
	session := cps.DBProvider.GetSession()

	stmt :=
		session.Query("DELETE FROM playlist_songs WHERE playlist_id = ? AND song_id = ?",
			playlistId, songId)

	if err := stmt.Exec(); err != nil {
		if err == gocql.ErrNotFound {
			return &customerrors.ErrNotFound{Message: err.Error()}
		}
		return &customerrors.ErrInternal{Message: err.Error()}
	}

	return nil
}

/*
// Returns first 100 playlist songs before a certain playlist song id.
// UUID validation is not provided
// May return Err internal
func (cps *CqlPlaylistSongRepository) GetPlaylistSongs(playlistId string, playlistSongId int64) ([]models.Song, error) {
	session := cps.DBProvider.GetSession()

	var songs []models.Song

	iter := session.Query(`
		SELECT user_id, chat_id, content, time_id, encrypted
		FROM playlist_songs
		WHERE playlist_id = ? AND id < ?
		LIMIT 100
		`, playlistId, playlistSongId).Iter()

	var (
		user_id   string
		chat_id   string
		content   string
		time_id   int64
		encrypted bool
	)

	for iter.Scan(&user_id, &chat_id, &content, &time_id, &encrypted) {
		songs = append(songs, models.Song{
			UserId:  user_id,
			ChatId:  chat_id,
			Content: content,
		})
	}

	if err := iter.Close(); err != nil {
		if err == gocql.ErrNotFound {
			return nil, &customerrors.ErrNotFound{Message: err.Error()}
		}
		return nil, &customerrors.ErrInternal{Message: err.Error()}
	}

	if len(songs) == 0 {
		return nil, &customerrors.ErrNotFound{Message: "no songs"}
	}

	return songs, nil
}
*/
