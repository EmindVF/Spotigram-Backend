package repositories

import (
	"spotigram/internal/customerrors"
	"spotigram/internal/infrastructure/abstractions"
	"spotigram/internal/service/models"
	"time"

	"github.com/gocql/gocql"
)

type CqlChatRepository struct {
	DBProvider abstractions.CqlDatabaseProvider
}

// Adds a message to the repository.
// May return ErrInternal on failure.
func (ccr *CqlChatRepository) AddMessage(m models.Message) error {
	session := ccr.DBProvider.GetSession()

	insertStmt := session.Query(`
		INSERT INTO messages (user_id, chat_id,creation_date, content, time_id, encrypted)
		VALUES (?, ?,?, ?, ?, ?)
	`)

	err := insertStmt.Bind(
		m.UserId,
		m.ChatId,
		m.Date,
		m.Content,
		m.TimeId,
		m.IsEncrypted,
	).Exec()

	if err != nil {
		return &customerrors.ErrInternal{Message: err.Error()}
	}

	return nil
}

// Gets at most 100 messages of the chat with offset.
// May return ErrInternal or ErrNotFound on failure.
func (ccr *CqlChatRepository) GetMessages(chatId string, timeId int64) ([]models.Message, error) {
	session := ccr.DBProvider.GetSession()

	var messages []models.Message

	iter := session.Query(`
		SELECT user_id, chat_id, creation_date ,content, time_id, encrypted 
		FROM messages 
		WHERE chat_id = ? AND time_id < ?
		LIMIT 100
		`, chatId, timeId).Iter()

	var (
		user_id       string
		chat_id       string
		creation_date time.Time
		content       string
		time_id       int64
		encrypted     bool
	)

	for iter.Scan(&user_id, &chat_id, &creation_date, &content, &time_id, &encrypted) {
		messages = append(messages, models.Message{
			UserId:      user_id,
			ChatId:      chat_id,
			Date:        creation_date,
			Content:     content,
			TimeId:      time_id,
			IsEncrypted: encrypted,
		})
	}

	if err := iter.Close(); err != nil {
		if err == gocql.ErrNotFound {
			return nil, &customerrors.ErrNotFound{Message: err.Error()}
		}
		return nil, &customerrors.ErrInternal{Message: err.Error()}
	}

	if len(messages) == 0 {
		return nil, &customerrors.ErrNotFound{Message: "no messages"}
	}

	return messages, nil
}

// Deletes a chat from the repository.
// May return ErrInternal or ErrNotFound on failure.
func (ccr *CqlChatRepository) DeleteChat(chatId string) error {
	session := ccr.DBProvider.GetSession()

	stmt :=
		session.Query("DELETE FROM messages WHERE chat_id = ?", chatId)

	if err := stmt.Exec(); err != nil {
		if err == gocql.ErrNotFound {
			return &customerrors.ErrNotFound{Message: err.Error()}
		}
		return &customerrors.ErrInternal{Message: err.Error()}
	}

	return nil
}

// Deletes a message from the repository.
// May return ErrInternal or ErrNotFound on failure.
func (ccr *CqlChatRepository) DeleteMessage(chatId string, messageId int64) error {
	session := ccr.DBProvider.GetSession()

	stmt :=
		session.Query("DELETE FROM messages WHERE chat_id = ? AND time_id = ?", chatId, messageId)

	if err := stmt.Exec(); err != nil {
		if err == gocql.ErrNotFound {
			return &customerrors.ErrNotFound{Message: err.Error()}
		}
		return &customerrors.ErrInternal{Message: err.Error()}
	}

	return nil
}
