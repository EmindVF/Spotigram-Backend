package abstractions

import "database/sql"

type DatabaseProvider interface {
	GetDb() *sql.DB
}
