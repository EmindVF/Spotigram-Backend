package repositories

import (
	"database/sql"
	"fmt"
	"spotigram/internal/customerrors"
	"spotigram/internal/infrastructure/abstractions"
	"spotigram/internal/service/models"
)

type SqlFriendRepository struct {
	DBProvider abstractions.SqlDatabaseProvider
}

// Returns every friend of a user by his id.
// May return ErrInternal or ErrNotFound on failure.
func (sdm *SqlFriendRepository) GetFriends(uuid1 string, offset int) ([]models.Friend, error) {
	var friends []models.Friend

	if offset < 0 {
		return nil, &customerrors.ErrInternal{Message: "invalid offset"}
	}

	db := sdm.DBProvider.GetDb()
	rows, err := db.Query("SELECT user2_id, chat_id FROM friendships WHERE user1_id = $1 OFFSET $2 LIMIT 100", uuid1, offset)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, &customerrors.ErrInternal{Message: err.Error()}
		} else {
			return nil, &customerrors.ErrNotFound{Message: "friends not found"}
		}
	}

	for rows.Next() {
		friends = append(friends, models.Friend{
			Id1: uuid1,
		})
		if err := rows.Scan(&friends[len(friends)-1].Id2, &friends[len(friends)-1].ChatId); err != nil {
			return nil, &customerrors.ErrInternal{Message: err.Error()}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, &customerrors.ErrInternal{Message: err.Error()}
	}

	rows.Close()

	rows, err = db.Query("SELECT user1_id, chat_id FROM friendships WHERE user2_id = $1", uuid1)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, &customerrors.ErrInternal{Message: err.Error()}
		}
	}

	for rows.Next() {
		friends = append(friends, models.Friend{
			Id1: uuid1,
		})
		if err := rows.Scan(&friends[len(friends)-1].Id2, &friends[len(friends)-1].ChatId); err != nil {
			return nil, &customerrors.ErrInternal{Message: err.Error()}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, &customerrors.ErrInternal{Message: err.Error()}
	}

	rows.Close()

	if len(friends) == 0 {
		return nil, &customerrors.ErrNotFound{Message: "friends not found"}
	}

	return friends, nil
}

// Returns a chat id of between friends.
// May return ErrInternal or ErrNotFound on failure.
func (sdm *SqlFriendRepository) GetChatIdByFriend(uuid1, uuid2 string) (chatId string, e error) {
	if uuid1 > uuid2 {
		uuid1, uuid2 = uuid2, uuid1
	}

	db := sdm.DBProvider.GetDb()
	row := db.QueryRow("SELECT chat_id FROM friendships WHERE user1_id = $1 AND user2_id = $2", uuid1, uuid2)
	if err := row.Scan(&chatId); err != nil {
		if err == sql.ErrNoRows {
			return "", &customerrors.ErrNotFound{Message: "no such friendship"}
		} else {
			return "", &customerrors.ErrInternal{Message: err.Error()}
		}
	}

	return chatId, nil
}

// Returns a friend by chat id.
// May return ErrInternal or ErrNotFound on failure.
func (sdm *SqlFriendRepository) GetFriendByChatId(uuid string) (f *models.Friend, e error) {
	f = &models.Friend{}
	db := sdm.DBProvider.GetDb()
	row := db.QueryRow("SELECT user1_id, user2_id, chat_id FROM friendships WHERE chat_id = $1", uuid)
	if err := row.Scan(&f.Id1, &f.Id2, &f.ChatId); err != nil {
		if err == sql.ErrNoRows {
			return nil, &customerrors.ErrNotFound{Message: "no such friendship"}
		} else {
			return nil, &customerrors.ErrInternal{Message: err.Error()}
		}
	}

	return f, nil
}

// Adds a friend to the repository.
// May return ErrInternal or ErrInvalidInput on failure.
func (sfm *SqlFriendRepository) AddFriend(f models.Friend) error {
	db := sfm.DBProvider.GetDb()

	stmt, err := db.Prepare("INSERT INTO friendships (user1_id, user2_id, chat_id) VALUES ($1, $2, $3)")
	if err != nil {
		panic(fmt.Errorf("error preparing AddFriend SQL statement: %v", err))
	}
	defer stmt.Close()

	if f.Id1 > f.Id2 {
		f.Id1, f.Id2 = f.Id2, f.Id1
	}

	_, err = stmt.Exec(f.Id1, f.Id2, f.ChatId)
	if err == sql.ErrConnDone {
		return &customerrors.ErrInternal{Message: "connection is done"}
	} else if err != nil {
		return &customerrors.ErrInvalidInput{Message: err.Error()}
	}

	return nil
}

// Deletes a friend from the repository.
// May return ErrInternal, ErrNotFound or ErrInvalidInput on failure.
func (sfm *SqlFriendRepository) DeleteFriend(uuid1, uuid2 string) error {
	db := sfm.DBProvider.GetDb()

	stmt, err := db.Prepare("DELETE FROM friendships WHERE user1_id = $1 AND user2_id = $2")
	if err != nil {
		panic(fmt.Errorf("error preparing DeleteFriend SQL statement: %v", err))
	}
	defer stmt.Close()

	if uuid1 > uuid2 {
		uuid1, uuid2 = uuid2, uuid1
	}

	res, err := stmt.Exec(uuid1, uuid2)
	if err == sql.ErrConnDone {
		return &customerrors.ErrInternal{Message: "connection is done"}
	} else if err != nil {
		return &customerrors.ErrInvalidInput{Message: err.Error()}
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return &customerrors.ErrInternal{Message: err.Error()}
	} else if rowsAffected < 1 {
		return &customerrors.ErrNotFound{Message: "friend request not found"}
	}

	return nil
}

// Checks whether a friend exits in the repository.
// May return ErrInternal on failure.
func (sfm *SqlFriendRepository) DoesFriendExist(uuid1, uuid2 string) (bool, error) {
	if uuid1 > uuid2 {
		uuid1, uuid2 = uuid2, uuid1
	}

	var result string
	err := sfm.DBProvider.GetDb().QueryRow(
		"SELECT user1_id FROM friendships WHERE user1_id = $1 AND user2_id = $2", uuid1, uuid2).Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, &customerrors.ErrInternal{Message: err.Error()}
		}
	}
	return true, nil
}
