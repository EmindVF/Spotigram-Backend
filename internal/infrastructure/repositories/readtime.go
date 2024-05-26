package repositories

import (
	"database/sql"
	"fmt"
	"spotigram/internal/customerrors"
	"spotigram/internal/infrastructure/abstractions"
	"spotigram/internal/service/models"
)

type SqlReadTimeRepository struct {
	DBProvider abstractions.SqlDatabaseProvider
}

func (r *SqlReadTimeRepository) GetReadTime(userId string, chatId string) (*models.ReadTime, error) {
	var timeId int64
	db := r.DBProvider.GetDb()
	row := db.QueryRow("SELECT time_id FROM read_times WHERE user_id = $1 AND chat_id = $2", userId, chatId)
	if err := row.Scan(&timeId); err != nil {
		if err == sql.ErrNoRows {
			return nil, &customerrors.ErrNotFound{Message: "no such readTime"}
		} else {
			return nil, &customerrors.ErrInternal{Message: err.Error()}
		}
	}

	return &models.ReadTime{
		UserId: userId,
		ChatId: chatId,
		TimeId: timeId,
	}, nil
}

func (r *SqlReadTimeRepository) DeleteReadTimeByChatId(chatId string) error {
	db := r.DBProvider.GetDb()

	stmt, err := db.Prepare("DELETE FROM read_times WHERE chat_id = $1")
	if err != nil {
		panic(fmt.Errorf("error preparing DeleteReadTimeByChatId SQL statement: %v", err))
	}
	defer stmt.Close()

	_, err = stmt.Exec(chatId)
	if err == sql.ErrConnDone {
		return &customerrors.ErrInternal{Message: "connection is done"}
	} else if err != nil {
		return &customerrors.ErrInvalidInput{Message: err.Error()}
	}
	return nil
}

func (r *SqlReadTimeRepository) AddReadTime(readTime *models.ReadTime) error {
	db := r.DBProvider.GetDb()

	stmt, err := db.Prepare(`
		INSERT INTO read_times (user_id, chat_id, time_id) 
		VALUES ($1, $2, $3)`)
	if err != nil {
		panic(fmt.Errorf("error preparing AddReadTime SQL statement: %v", err))
	}
	defer stmt.Close()

	_, err = stmt.Exec(readTime.UserId, readTime.ChatId, readTime.TimeId)
	if err == sql.ErrConnDone {
		return &customerrors.ErrInternal{Message: "connection is done"}
	} else if err != nil {
		return &customerrors.ErrInvalidInput{Message: err.Error()}
	}

	return nil
}

func (r *SqlReadTimeRepository) UpdateReadTime(userId string, chatId string, timeId int64) error {
	db := r.DBProvider.GetDb()

	stmt, err := db.Prepare("UPDATE read_times SET time_id = $1 WHERE user_id = $2 AND chat_id = $3")
	if err != nil {
		panic(fmt.Errorf("error preparing UpdateReadTime SQL statement: %v", err))
	}
	defer stmt.Close()

	_, err = stmt.Exec(timeId, userId, chatId)
	if err == sql.ErrConnDone {
		return &customerrors.ErrInternal{Message: "connection is done"}
	} else if err != nil {
		return &customerrors.ErrInvalidInput{Message: err.Error()}
	}

	return nil
}
