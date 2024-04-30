package abstractions

import "github.com/gocql/gocql"

type CqlDatabaseProvider interface {
	GetSession() *gocql.Session
}
