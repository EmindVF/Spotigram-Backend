package repositories

import (
	"database/sql"
	"fmt"
	"spotigram/internal/customerrors"
	"spotigram/internal/infrastructure/abstractions"
	"spotigram/internal/service/models"
)

type SqlFriendRequestRepository struct {
	DBProvider abstractions.SqlDatabaseProvider
}

// Returns every friend request sent by a user.
// May return ErrInternal or ErrNotFound on failure.
func (sfrr *SqlFriendRequestRepository) GetFriendRequestsSent(uuid1 string, offset int) ([]models.FriendRequest, error) {
	if offset < 0 {
		return nil, &customerrors.ErrInternal{Message: "invalid offset"}
	}
	var friendRequests []models.FriendRequest
	db := sfrr.DBProvider.GetDb()
	rows, err := db.Query("SELECT recipient_id, is_ignored FROM friend_requests WHERE sender_id = $1 LIMIT 100 OFFSET $2", uuid1, offset)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, &customerrors.ErrInternal{Message: err.Error()}
		} else {
			return nil, &customerrors.ErrNotFound{Message: "friend requests not found"}
		}
	}

	for rows.Next() {
		friendRequests = append(friendRequests, models.FriendRequest{
			SenderId: uuid1,
		})
		if err := rows.Scan(&friendRequests[len(friendRequests)-1].RecipientId,
			&friendRequests[len(friendRequests)-1].IsIgnored); err != nil {
			return nil, &customerrors.ErrInternal{Message: err.Error()}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, &customerrors.ErrInternal{Message: err.Error()}
	}

	rows.Close()

	if len(friendRequests) == 0 {
		return nil, &customerrors.ErrNotFound{Message: "friend requests not found"}
	}

	return friendRequests, nil
}

// Returns every friend request recieved by a user.
// May return ErrInternal or ErrNotFound on failure.
func (sfrr *SqlFriendRequestRepository) GetFriendRequestsReceived(uuid1 string, offset int) ([]models.FriendRequest, error) {
	if offset < 0 {
		return nil, &customerrors.ErrInternal{Message: "invalid offset"}
	}
	var friendRequests []models.FriendRequest
	db := sfrr.DBProvider.GetDb()
	rows, err := db.Query("SELECT sender_id, is_ignored FROM friend_requests WHERE recipient_id = $1 LIMIT 100 OFFSET $2", uuid1, offset)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, &customerrors.ErrInternal{Message: err.Error()}
		} else {
			return nil, &customerrors.ErrNotFound{Message: "friends requests not found"}
		}
	}

	for rows.Next() {
		friendRequests = append(friendRequests, models.FriendRequest{
			RecipientId: uuid1,
		})
		if err := rows.Scan(&friendRequests[len(friendRequests)-1].SenderId,
			&friendRequests[len(friendRequests)-1].IsIgnored); err != nil {
			return nil, &customerrors.ErrInternal{Message: err.Error()}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, &customerrors.ErrInternal{Message: err.Error()}
	}

	rows.Close()

	if len(friendRequests) == 0 {
		return nil, &customerrors.ErrNotFound{Message: "friend requests not found"}
	}

	return friendRequests, nil
}

// Adds a friend request to the repository.
// May return ErrInternal or ErrInvalidInput on failure.
func (sfm *SqlFriendRequestRepository) AddFriendRequest(f models.FriendRequest) error {
	db := sfm.DBProvider.GetDb()

	stmt, err := db.Prepare("INSERT INTO friend_requests (sender_id, recipient_id, is_ignored) VALUES ($1, $2, $3)")
	if err != nil {
		panic(fmt.Errorf("error preparing AddFriendRequest SQL statement: %v", err))
	}
	defer stmt.Close()

	_, err = stmt.Exec(f.SenderId, f.RecipientId, f.IsIgnored)
	if err == sql.ErrConnDone {
		return &customerrors.ErrInternal{Message: "connection is done"}
	} else if err != nil {
		return &customerrors.ErrInvalidInput{Message: err.Error()}
	}

	return nil
}

// Updates a friend request's is_ignored.
// UUID validation is not provided.
// May return ErrInternal or ErrNotFound on failure.
func (sdm *SqlFriendRequestRepository) UpdateIsIgnored(sender, recipient string, isIgnored bool) error {
	res, err := sdm.DBProvider.GetDb().Exec(
		"UPDATE friend_requests SET is_ignored = $1 WHERE sender_id = $2 AND recipient_id = $3",
		isIgnored, sender, recipient)
	if err != nil {
		return &customerrors.ErrInternal{Message: err.Error()}
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return &customerrors.ErrInternal{Message: err.Error()}
	} else if rowsAffected < 1 {
		return &customerrors.ErrNotFound{Message: "friend request not found"}
	}
	return nil
}

// Deletes a friend request from the repository.
// May return ErrInternal, ErrNotFound or ErrInvalidInput on failure.
func (sfm *SqlFriendRequestRepository) DeleteFriendRequest(sender, recipient string) error {
	db := sfm.DBProvider.GetDb()

	stmt, err := db.Prepare("DELETE FROM friend_requests WHERE sender_id = $1 AND recipient_id = $2")
	if err != nil {
		panic(fmt.Errorf("error preparing DeleteFriendRequest SQL statement: %v", err))
	}
	defer stmt.Close()

	res, err := stmt.Exec(sender, recipient)
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

// Checks whether a friend request exits in the repository.
// May return ErrInternal on failure.
func (sfm *SqlFriendRequestRepository) DoesFriendRequestExist(sender, recipient string) (bool, error) {
	var result string
	err := sfm.DBProvider.GetDb().QueryRow(
		"SELECT sender_id FROM friend_requests WHERE sender_id = $1 AND recipient_id = $2",
		sender, recipient).Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, &customerrors.ErrInternal{Message: err.Error()}
		}
	}
	return true, nil
}
