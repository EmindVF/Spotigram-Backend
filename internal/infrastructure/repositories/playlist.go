package repositories

import (
	"database/sql"
	"fmt"
	"spotigram/internal/customerrors"
	"spotigram/internal/infrastructure/abstractions"
	"spotigram/internal/service/models"
)

type SqlPlaylistRepository struct {
	DBProvider abstractions.SqlDatabaseProvider
}

// Returns first 100 messages before a certain time id.
// UUID validation is not provided
// May return ErrInternal
func (spp *SqlPlaylistRepository) GetPlaylists(userId string, offset int) ([]models.Playlist, error) {
	var playlists []models.Playlist

	if offset < 0 {
		return nil, &customerrors.ErrInternal{Message: "invalid offset"}
	}

	db := spp.DBProvider.GetDb()
	rows, err := db.Query("SELECT id, name FROM playlists WHERE user_id = $1 OFFSET $2 LIMIT 100", userId, offset)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, &customerrors.ErrInternal{Message: err.Error()}
		} else {
			return nil, &customerrors.ErrNotFound{Message: "playlists not found"}
		}
	}

	for rows.Next() {
		playlists = append(playlists, models.Playlist{})
		if err := rows.Scan(
			&playlists[len(playlists)-1].Id,
			&playlists[len(playlists)-1].Name); err != nil {
			return nil, &customerrors.ErrInternal{Message: err.Error()}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, &customerrors.ErrInternal{Message: err.Error()}
	}

	rows.Close()

	if len(playlists) == 0 {
		return nil, &customerrors.ErrNotFound{Message: "playlists not found"}
	}

	return playlists, nil
}

// Returns the user id of playlist creator by its id.
// UUID validation is not provided
func (spp *SqlPlaylistRepository) GetUserIdByPlaylist(playlistId string) (string, error) {
	var userId string
	db := spp.DBProvider.GetDb()
	row := db.QueryRow("SELECT user_id FROM playlists WHERE id = $1", playlistId)
	if err := row.Scan(&userId); err != nil {
		if err == sql.ErrNoRows {
			return "", &customerrors.ErrNotFound{Message: "no such playlist"}
		} else {
			return "", &customerrors.ErrInternal{Message: err.Error()}
		}
	}

	return userId, nil
}

// Adds a playlist.
// May return ErrInternal or ErrInvalidInput
func (spp *SqlPlaylistRepository) AddPlaylist(p models.Playlist) error {
	db := spp.DBProvider.GetDb()

	stmt, err := db.Prepare("INSERT INTO playlists (id, user_id, name) VALUES ($1, $2, $3)")
	if err != nil {
		panic(fmt.Errorf("error preparing AddFriend SQL statement: %v", err))
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.Id, p.UserId, p.Name)
	if err == sql.ErrConnDone {
		return &customerrors.ErrInternal{Message: "connection is done"}
	} else if err != nil {
		return &customerrors.ErrInvalidInput{Message: err.Error()}
	}

	return nil
}

// Deletes a playlist.
// May return ErrInternal or ErrNotFound
func (spp *SqlPlaylistRepository) DeletePlaylist(playlistId string) error {
	db := spp.DBProvider.GetDb()

	stmt, err := db.Prepare("DELETE FROM playlists WHERE id = $1")
	if err != nil {
		panic(fmt.Errorf("error preparing DeleteFriend SQL statement: %v", err))
	}
	defer stmt.Close()

	res, err := stmt.Exec(playlistId)
	if err == sql.ErrConnDone {
		return &customerrors.ErrInternal{Message: "connection is done"}
	} else if err != nil {
		return &customerrors.ErrInvalidInput{Message: err.Error()}
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return &customerrors.ErrInternal{Message: err.Error()}
	} else if rowsAffected < 1 {
		return &customerrors.ErrNotFound{Message: "playlist not found"}
	}

	return nil
}
