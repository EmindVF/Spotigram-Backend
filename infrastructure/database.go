package infrastructure

import "database/sql"

var (
	Database DatabaseProvider
)

type DatabaseProvider interface {
	GetDb() *sql.DB
}
