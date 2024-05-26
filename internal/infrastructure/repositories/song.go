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

func (ssr *SqlSongRepository) GetSongs(offset int, songNameFilter string, creatorIdFilter string) ([]models.Song, error) {
	if offset < 0 {
		return nil, &customerrors.ErrInternal{Message: "invalid offset"}
	}

	var songs []models.Song
	db := ssr.DBProvider.GetDb()
	var rows *sql.Rows
	var err error
	if creatorIdFilter == "" {
		rows, err = db.Query(
			"SELECT id, creator_id, name, length, streams FROM songs WHERE name LIKE '%' || $2 || '%' ORDER BY streams DESC OFFSET $1 LIMIT 100",
			offset, songNameFilter)
	} else {
		rows, err = db.Query(
			"SELECT id, creator_id, name, length, streams FROM songs WHERE name LIKE '%' || $2 || '%' AND creator_id = $3 ORDER BY streams DESC OFFSET $1 LIMIT 100",
			offset, songNameFilter, creatorIdFilter)
	}
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
			&songs[len(songs)-1].Name,
			&songs[len(songs)-1].Length,
			&songs[len(songs)-1].Streams); err != nil {
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

func (ssr *SqlSongRepository) GetSongInfo(songId string) (*models.Song, error) {
	song := models.Song{
		Id: songId,
	}
	db := ssr.DBProvider.GetDb()
	row := db.QueryRow("SELECT name, length, creator_id, streams FROM songs WHERE id = $1", songId)
	if err := row.Scan(&song.Name, &song.Length, &song.CreatorId, &song.Streams); err != nil {
		if err == sql.ErrNoRows {
			return nil, &customerrors.ErrNotFound{Message: "no such song"}
		} else {
			return nil, &customerrors.ErrInternal{Message: err.Error()}
		}
	}

	return &song, nil
}

func (ssr *SqlSongRepository) GetSongFile(songId string) ([]byte, error) {
	var file []byte
	db := ssr.DBProvider.GetDb()
	row := db.QueryRow("SELECT file FROM songs WHERE id = $1", songId)
	if err := row.Scan(&file); err != nil {
		if err == sql.ErrNoRows {
			return nil, &customerrors.ErrNotFound{Message: "no such song"}
		} else {
			return nil, &customerrors.ErrInternal{Message: err.Error()}
		}
	}

	return file, nil
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

func (ssr *SqlSongRepository) AddSong(song models.Song, picture []byte, file []byte) error {
	db := ssr.DBProvider.GetDb()

	stmt, err := db.Prepare("INSERT INTO songs (id, creator_id, name, length, picture, file, streams) VALUES ($1, $2, $3, $4, $5, $6, 0)")
	if err != nil {
		panic(fmt.Errorf("error preparing AddSong SQL statement: %v", err))
	}
	defer stmt.Close()

	_, err = stmt.Exec(song.Id, song.CreatorId, song.Name, song.Length, picture, file)
	if err == sql.ErrConnDone {
		return &customerrors.ErrInternal{Message: "connection is done"}
	} else if err != nil {
		return &customerrors.ErrInvalidInput{Message: err.Error()}
	}

	return nil
}

// Adds a playlist.
// May return ErrInternal or ErrInvalidInput
func (spp *SqlSongRepository) UpdateSongName(songId string, name string) error {
	db := spp.DBProvider.GetDb()

	stmt, err := db.Prepare("UPDATE songs SET name = $1 WHERE id = $2")
	if err != nil {
		panic(fmt.Errorf("error preparing UpdateSongName SQL statement: %v", err))
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, songId)
	if err == sql.ErrConnDone {
		return &customerrors.ErrInternal{Message: "connection is done"}
	} else if err != nil {
		return &customerrors.ErrInvalidInput{Message: err.Error()}
	}

	return nil
}

func (ssr *SqlSongRepository) IncrementStreams(songId string) error {
	db := ssr.DBProvider.GetDb()

	stmt, err := db.Prepare("UPDATE songs SET streams = streams + 1 WHERE id = $1")
	if err != nil {
		panic(fmt.Errorf("error preparing IncrementStreams SQL statement: %v", err))
	}
	defer stmt.Close()

	_, err = stmt.Exec(songId)
	if err == sql.ErrConnDone {
		return &customerrors.ErrInternal{Message: "connection is done"}
	} else if err != nil {
		return &customerrors.ErrNotFound{Message: err.Error()}
	}

	return nil
}

// Returns a song's picture by its uuid.
// UUID validation is not provided.
// May return ErrInternal or ErrNotFound on failure.
func (sdm *SqlSongRepository) GetPicture(songId string) ([]byte, error) {
	var pic []byte
	err := sdm.DBProvider.GetDb().QueryRow(
		"SELECT picture FROM songs WHERE id = $1", songId).Scan(
		&pic)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &customerrors.ErrNotFound{Message: "song not found"}
		} else {
			return nil, &customerrors.ErrInternal{Message: err.Error()}
		}
	}
	return pic, nil
}
