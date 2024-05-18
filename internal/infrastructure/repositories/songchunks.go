package repositories

import (
	"spotigram/internal/customerrors"
	"spotigram/internal/infrastructure/abstractions"

	"github.com/gocql/gocql"
)

type CqlSongChunkRepository struct {
	DBProvider abstractions.CqlDatabaseProvider
}

// Deletes song chunks from the repository.
// May return ErrInternal or ErrNotFound on failure.
func (cps *CqlSongChunkRepository) DeleteSongChunks(songId string) error {
	session := cps.DBProvider.GetSession()

	stmt :=
		session.Query("DELETE FROM song_chunks WHERE song_id = ?",
			songId)

	if err := stmt.Exec(); err != nil {
		if err == gocql.ErrNotFound {
			return &customerrors.ErrNotFound{Message: err.Error()}
		}
		return &customerrors.ErrInternal{Message: err.Error()}
	}

	return nil
}

// Deletes song chunks from the repository.
// May return ErrInternal or ErrNotFound on failure.
func (cps *CqlSongChunkRepository) GetSongChunk(songId string, index int) ([]byte, error) {
	session := cps.DBProvider.GetSession()
	var file []byte
	err := session.Query("SELECT file FROM song_chunks WHERE song_id = ? AND ind = ?",
		songId, index).Scan(&file)
	if err != nil {
		if err == gocql.ErrNotFound {
			return nil, &customerrors.ErrNotFound{Message: err.Error()}
		}
		return nil, &customerrors.ErrInternal{Message: err.Error()}
	}

	return file, nil
}

// Deletes a song from the repository.
// May return ErrInternal or ErrNotFound on failure.
func (cps *CqlSongChunkRepository) AddSongChunk(songId string, ind int, chunk []byte) error {
	session := cps.DBProvider.GetSession()

	insertStmt := session.Query(`
		INSERT INTO song_chunks (song_id, ind, file)
		VALUES (?, ?, ?)
	`)
	err := insertStmt.Bind(
		songId, ind, chunk,
	).Exec()

	if err != nil {
		return &customerrors.ErrInternal{Message: err.Error()}
	}

	return nil
}
