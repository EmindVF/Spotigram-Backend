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
		INSERT INTO playlist_songs (playlist_id, song_id)
		VALUES (?, ?)
	`)

	err := insertStmt.Bind(
		playlistSong.PlaylistId,
		playlistSong.SongId,
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
func (cps *CqlPlaylistSongRepository) DeleteSong(songId string) error {
	session := cps.DBProvider.GetSession()

	stmt :=
		session.Query("DELETE FROM playlist_songs WHERE song_id = ?",
			songId)

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
func (cps *CqlPlaylistSongRepository) GetPlaylistLength(playlistId string) (int, error) {
	session := cps.DBProvider.GetSession()

	var count int
	query := session.Query("SELECT COUNT(*) FROM playlist_songs WHERE playlist_id = ?", playlistId)
	err := query.Scan(&count)
	if err != nil {
		if err == gocql.ErrNotFound {
			return 0, &customerrors.ErrNotFound{Message: err.Error()}
		}
		return 0, &customerrors.ErrInternal{Message: err.Error()}
	}
	return count, nil

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
func (cps *CqlPlaylistSongRepository) GetPlaylistSongs(playlistId string) ([]models.PlaylistSong, error) {
	session := cps.DBProvider.GetSession()

	var songIds []models.PlaylistSong

	iter := session.Query(`
		SELECT song_id
		FROM playlist_songs
		WHERE playlist_id = ? 
		LIMIT 100
		`, playlistId).Iter()

	var (
		song_id string
	)

	for iter.Scan(&song_id) {
		songIds = append(songIds, models.PlaylistSong{
			PlaylistId: playlistId,
			SongId:     song_id,
		})
	}

	if err := iter.Close(); err != nil {
		if err == gocql.ErrNotFound {
			return nil, &customerrors.ErrNotFound{Message: err.Error()}
		}
		return nil, &customerrors.ErrInternal{Message: err.Error()}
	}

	if len(songIds) == 0 {
		return nil, &customerrors.ErrNotFound{Message: "no songs"}
	}

	return songIds, nil
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
