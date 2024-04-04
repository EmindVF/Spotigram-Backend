package repositories

import (
	"database/sql"
	"fmt"
	"spotigram/internal/customerrors"
	"spotigram/internal/infrastructure/abstractions"
	"spotigram/internal/service/models"
)

type SqlUserRepository struct {
	DBProvider abstractions.DatabaseProvider
}

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

func (sdm *SqlUserRepository) GetUser(uuid string) (*models.User, error) {
	var user = models.User{}
	err := sdm.DBProvider.GetDb().QueryRow(
		"SELECT id, name, email, password, verified FROM users WHERE id = $1", uuid).Scan(
		&user.Id, &user.Name, &user.Email, &user.Password, &user.Verified)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &customerrors.ErrNotFound{Message: "no such user"}
		} else {
			return nil, &customerrors.ErrInternal{Message: err.Error()}
		}
	}
	return &user, nil
}
