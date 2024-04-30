package abstractions

import "database/sql"

type SqlDatabaseProvider interface {
	GetDb() *sql.DB
}
