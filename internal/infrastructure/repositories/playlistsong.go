package repositories

import (
	"spotigram/internal/customerrors"
	"spotigram/internal/infrastructure/abstractions"
	"spotigram/internal/service/models"

	"github.com/gocql/gocql"
)

type CqlPlaylistSongRepository struct {
	DBProvider abstractions.CqlDatabaseProvider
}

// Adds a song to the playlist.
// May return ErrInternal on failure.
func (cps *CqlPlaylistSongRepository) AddPlaylistSong(playlistSong models.PlaylistSong) error {
	session := cps.DBProvider.GetSession()

	insertStmt := session.Query(`
		INSERT INTO playlist_songs (playlist_id, song_id, user_id, name, length)
		VALUES (?, ?, ?, ?)
	`)

	err := insertStmt.Bind(
		playlistSong.PlaylistId,
		playlistSong.SongId,
		playlistSong.UserId,
		playlistSong.Name,
		playlistSong.Length,
	).Exec()

	if err != nil {
		return &customerrors.ErrInternal{Message: err.Error()}
	}

	return nil
}

// Deletes a song from the repository.
// May return ErrInternal or ErrNotFound on failure.
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

// Deletes a song from the repository.
// May return ErrInternal or ErrNotFound on failure.
func (cps *CqlPlaylistSongRepository) DeletePlaylistSongs(playlistId string) error {
	session := cps.DBProvider.GetSession()

	stmt :=
		session.Query("DELETE FROM playlist_songs WHERE playlist_id = ?",
			playlistId)

	if err := stmt.Exec(); err != nil {
		if err == gocql.ErrNotFound {
			return &customerrors.ErrNotFound{Message: err.Error()}
		}
		return &customerrors.ErrInternal{Message: err.Error()}
	}

	return nil
}

// Returns first 100 playlist songs.
// UUID validation is not provided
// May return Err internal
func (cps *CqlPlaylistSongRepository) GetPlaylistSongs(playlistId string) ([]models.Song, error) {
	session := cps.DBProvider.GetSession()

	var songs []models.Song

	iter := session.Query(`
		SELECT song_id, user_id, name, length
		FROM playlist_songs
		WHERE playlist_id = ? 
		LIMIT 100
		`, playlistId).Iter()

	var (
		song_id string
		user_id string
		name    string
		length  int
	)

	for iter.Scan(&song_id, &user_id, &name, &length) {
		songs = append(songs, models.Song{
			Id:        song_id,
			CreatorId: user_id,
			Name:      name,
			Length:    length,
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

// Checks if the song is in a playlist.
// May return ErrInternal or ErrNotFound on failure.
func (ccr *CqlPlaylistSongRepository) IsSongInPlaylist(playlistId string, songId string) (bool, error) {
	session := ccr.DBProvider.GetSession()

	iter := session.Query(`
		SELECT playlist_id
		FROM playlist_songs
		WHERE playlist_id = ? AND song_id = ?
		LIMIT 1
		`, playlistId, songId).Iter()

	var newPlaylistId string
	var songs []string

	for iter.Scan(&newPlaylistId) {
		songs = append(songs, newPlaylistId)
	}

	if err := iter.Close(); err != nil {
		if err == gocql.ErrNotFound {
			return false, nil
		}
		return false, &customerrors.ErrInternal{Message: err.Error()}
	}

	if len(songs) == 0 {
		return false, nil
	}

	return true, nil
}
