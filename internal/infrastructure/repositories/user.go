package repositories

import (
	"database/sql"
	"fmt"
	"spotigram/internal/customerrors"
	"spotigram/internal/infrastructure/abstractions"
	"spotigram/internal/service/models"
)

type SqlUserRepository struct {
	DBProvider abstractions.SqlDatabaseProvider
}

// Adds user to the repository.
// May return ErrInternal or ErrInvalidInput on failure.
func (sdm *SqlUserRepository) AddUser(sud models.User) error {
	db := sdm.DBProvider.GetDb()

	stmt, err := db.Prepare("INSERT INTO users (id, name, email, password, picture, verified, public_key) VALUES ($1, $2, $3, $4, $5, $6, $7)")
	if err != nil {
		panic(fmt.Errorf("error preparing AddUser SQL statement: %v", err))
	}
	defer stmt.Close()

	_, err = stmt.Exec(sud.Id, sud.Name, sud.Email, sud.Password, nil, sud.Verified, nil)
	if err == sql.ErrConnDone {
		return &customerrors.ErrInternal{Message: "connection is done"}
	} else if err != nil {
		return &customerrors.ErrInvalidInput{Message: err.Error()}
	}

	return nil
}

// Returns uuid and hashed password of an user by its email.
// Email validation is not provided.
// May return ErrInternal or ErrNotFound on failure.
func (sdm *SqlUserRepository) GetUUIDAndPasswordByEmail(email string) (uuid string, passwordHash string, e error) {
	db := sdm.DBProvider.GetDb()
	row := db.QueryRow("SELECT id, password FROM users WHERE email = $1", email)
	if err := row.Scan(&uuid, &passwordHash); err != nil {
		if err == sql.ErrNoRows {
			return "", "", &customerrors.ErrNotFound{Message: "no such email"}
		} else {
			return "", "", &customerrors.ErrInternal{Message: err.Error()}
		}
	}

	return uuid, passwordHash, nil
}

// Returns a user by its uuid.
// UUID validation is not provided.
// May return ErrInternal or ErrNotFound on failure.
func (sdm *SqlUserRepository) GetPassword(uuid string) (string, error) {
	var p string
	err := sdm.DBProvider.GetDb().QueryRow(
		"SELECT password FROM users WHERE id = $1", uuid).Scan(&p)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", &customerrors.ErrNotFound{Message: "user not found"}
		} else {
			return "", &customerrors.ErrInternal{Message: err.Error()}
		}
	}
	return p, nil
}

// Returns a user's public key by its uuid.
// UUID validation is not provided.
// May return ErrInternal or ErrNotFound on failure.
func (sdm *SqlUserRepository) GetPublicKey(uuid string) (string, error) {
	var p string
	err := sdm.DBProvider.GetDb().QueryRow(
		"SELECT public_key FROM users WHERE id = $1", uuid).Scan(&p)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", &customerrors.ErrNotFound{Message: "user not found"}
		} else {
			return "", &customerrors.ErrInternal{Message: err.Error()}
		}
	}
	return p, nil
}

// Returns a user by its uuid.
// UUID validation is not provided.
// May return ErrInternal or ErrNotFound on failure.
func (sdm *SqlUserRepository) GetUser(uuid string) (*models.User, error) {
	var user = models.User{}
	err := sdm.DBProvider.GetDb().QueryRow(
		"SELECT id, name, email, password, verified FROM users WHERE id = $1", uuid).Scan(
		&user.Id, &user.Name, &user.Email, &user.Password, &user.Verified)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &customerrors.ErrNotFound{Message: "user not found"}
		} else {
			return nil, &customerrors.ErrInternal{Message: err.Error()}
		}
	}

	return &user, nil
}

// Returns a users list, given offset and filter.
// UUID validation is not provided.
// May return ErrInternal or ErrNotFound on failure.
func (sdm *SqlUserRepository) GetUsers(offset int, usernameFilter string) ([]models.User, error) {
	if offset < 0 {
		return nil, &customerrors.ErrInternal{Message: "invalid offset"}
	}

	var users []models.User
	db := sdm.DBProvider.GetDb()
	rows, err := db.Query(
		"SELECT id, name, email, password, verified FROM users WHERE name LIKE '%' || $2 || '%' OFFSET $1 LIMIT 100",
		offset, usernameFilter)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &customerrors.ErrNotFound{Message: "users not found"}
		} else {
			return nil, &customerrors.ErrInternal{Message: err.Error()}
		}
	}

	for rows.Next() {
		users = append(users, models.User{})
		if err := rows.Scan(
			&users[len(users)-1].Id,
			&users[len(users)-1].Name,
			&users[len(users)-1].Email,
			&users[len(users)-1].Password,
			&users[len(users)-1].Verified); err != nil {
			return nil, &customerrors.ErrInternal{Message: err.Error()}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, &customerrors.ErrInternal{Message: err.Error()}
	}

	rows.Close()

	if len(users) == 0 {
		return nil, &customerrors.ErrNotFound{Message: "users not found"}
	}

	return users, nil
}

// Returns a user's picture by its uuid.
// UUID validation is not provided.
// May return ErrInternal or ErrNotFound on failure.
func (sdm *SqlUserRepository) GetPicture(uuid string) ([]byte, error) {
	var pic []byte
	err := sdm.DBProvider.GetDb().QueryRow(
		"SELECT picture FROM users WHERE id = $1", uuid).Scan(
		&pic)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &customerrors.ErrNotFound{Message: "user not found"}
		} else {
			return nil, &customerrors.ErrInternal{Message: err.Error()}
		}
	}
	return pic, nil
}

// Updates a user's name by its uuid.
// UUID validation is not provided.
// May return ErrInternal or ErrNotFound on failure.
func (sdm *SqlUserRepository) UpdateName(uuid string, name string) error {
	res, err := sdm.DBProvider.GetDb().Exec(
		"UPDATE users SET name = $1 WHERE id = $2", name, uuid)
	if err != nil {
		return &customerrors.ErrInternal{Message: err.Error()}
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return &customerrors.ErrInternal{Message: err.Error()}
	} else if rowsAffected < 1 {
		return &customerrors.ErrNotFound{Message: "user not found"}
	}
	return nil
}

// Updates a user's password by its uuid.
// UUID validation is not provided.
// May return ErrInternal or ErrNotFound on failure.
func (sdm *SqlUserRepository) UpdatePassword(uuid string, password string) error {
	res, err := sdm.DBProvider.GetDb().Exec(
		"UPDATE users SET password = $1 WHERE id = $2", password, uuid)
	if err != nil {
		return &customerrors.ErrInternal{Message: err.Error()}
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return &customerrors.ErrInternal{Message: err.Error()}
	} else if rowsAffected < 1 {
		return &customerrors.ErrNotFound{Message: "user not found"}
	}
	return nil
}

// Updates a user's public_key by its uuid.
// UUID validation is not provided, public_key validation is not provided.
// May return ErrInternal or ErrNotFound on failure.
func (sdm *SqlUserRepository) UpdatePublicKey(uuid string, public_key string) error {
	res, err := sdm.DBProvider.GetDb().Exec(
		"UPDATE users SET public_key = $1 WHERE id = $2", public_key, uuid)
	if err != nil {
		return &customerrors.ErrInternal{Message: err.Error()}
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return &customerrors.ErrInternal{Message: err.Error()}
	} else if rowsAffected < 1 {
		return &customerrors.ErrNotFound{Message: "user not found"}
	}
	return nil
}

// Updates a user's picture by its uuid.
// UUID validation is not provided, image validation is not provided.
// May return ErrInternal or ErrNotFound on failure.
func (sdm *SqlUserRepository) UpdatePicture(uuid string, image []byte) error {
	res, err := sdm.DBProvider.GetDb().Exec(
		"UPDATE users SET picture = $1 WHERE id = $2", image, uuid)
	if err != nil {
		return &customerrors.ErrInternal{Message: err.Error()}
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return &customerrors.ErrInternal{Message: err.Error()}
	} else if rowsAffected < 1 {
		return &customerrors.ErrNotFound{Message: "user not found"}
	}
	return nil
}

// Returns bool on whether the user uuid is present.
// UUID validation is not provided.
// May return ErrInternal on failure.
func (sdm *SqlUserRepository) DoesUserExist(uuid string) (bool, error) {
	var result string
	err := sdm.DBProvider.GetDb().QueryRow(
		"SELECT id FROM users WHERE id = $1", uuid).Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, &customerrors.ErrInternal{Message: err.Error()}
		}
	}
	return true, nil
}
