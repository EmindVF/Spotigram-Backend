package repositories

import (
	"database/sql"
	"fmt"
	"spotigram/internal/customerrors"
	"spotigram/internal/infrastructure/abstractions"
	"spotigram/internal/service/models"
)

type SqlSongRepository struct {
	DBProvider abstractions.SqlDatabaseProvider
}

func (ssr *SqlSongRepository) GetSongs(offset int) ([]models.Song, error) {
	if offset < 0 {
		return nil, &customerrors.ErrInternal{Message: "invalid offset"}
	}

	var songs []models.Song
	db := ssr.DBProvider.GetDb()
	rows, err := db.Query("SELECT id, creator_id, name FROM songs OFFSET $1 LIMIT 100", offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &customerrors.ErrNotFound{Message: "songs not found"}
		} else {
			return nil, &customerrors.ErrInternal{Message: err.Error()}
		}
	}

	for rows.Next() {
		songs = append(songs, models.Song{})
		if err := rows.Scan(
			&songs[len(songs)-1].Id,
			&songs[len(songs)-1].CreatorId,
			&songs[len(songs)-1].Name); err != nil {
			return nil, &customerrors.ErrInternal{Message: err.Error()}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, &customerrors.ErrInternal{Message: err.Error()}
	}

	rows.Close()

	if len(songs) == 0 {
		return nil, &customerrors.ErrNotFound{Message: "songs not found"}
	}

	return songs, nil
}

func (ssr *SqlSongRepository) DeleteSong(songId string) error {
	db := ssr.DBProvider.GetDb()

	stmt, err := db.Prepare("DELETE FROM songs WHERE id = $1")
	if err != nil {
		panic(fmt.Errorf("error preparing DeleteSong SQL statement: %v", err))
	}
	defer stmt.Close()

	res, err := stmt.Exec(songId)
	if err == sql.ErrConnDone {
		return &customerrors.ErrInternal{Message: "connection is done"}
	} else if err != nil {
		return &customerrors.ErrInvalidInput{Message: err.Error()}
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return &customerrors.ErrInternal{Message: err.Error()}
	} else if rowsAffected < 1 {
		return &customerrors.ErrNotFound{Message: "song not found"}
	}

	return nil
}

func (ssr *SqlSongRepository) AddSong(song models.Song) error {
	db := ssr.DBProvider.GetDb()

	stmt, err := db.Prepare("INSERT INTO songs (id, creator_id, name, file) VALUES ($1, $2, $3, $4)")
	if err != nil {
		panic(fmt.Errorf("error preparing AddSong SQL statement: %v", err))
	}
	defer stmt.Close()

	_, err = stmt.Exec(song.Id, song.CreatorId, song.Name, song.File)
	if err == sql.ErrConnDone {
		return &customerrors.ErrInternal{Message: "connection is done"}
	} else if err != nil {
		return &customerrors.ErrInvalidInput{Message: err.Error()}
	}

	return nil
}
