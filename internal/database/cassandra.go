package database

import (
	"fmt"
	"os"
	"spotigram/internal/config"
	"strings"

	"github.com/gocql/gocql"
)

type cassandraCqlDatabaseProvider struct {
	session *gocql.Session
}

func (p *cassandraCqlDatabaseProvider) GetSession() *gocql.Session {
	return p.session
}

// Initializes connected cassandra database by the parameters in the config.
func initializeCassandraKeyspace(cfg *config.Config, s *gocql.Session) {
	script, err := os.ReadFile(cfg.CqlDb.InitKeyspaceScriptPath)
	if err != nil {
		panic(fmt.Errorf("fatal error reading initialize cassandra cql keyspace script: %v", err))
	}

	statements := strings.Split(string(script), ";")

	for _, statement := range statements {
		trimmedStatement := strings.TrimSpace(statement)
		if trimmedStatement != "" {
			err := s.Query(trimmedStatement).Exec()
			if err != nil {
				panic(fmt.Errorf("fatal error executing initialize cassandra cql keyspace script statement: %v", err))
			}
		}
	}
}

// Initializes connected cassandra database tables.
func initializeCassandraTable(cfg *config.Config, s *gocql.Session) {
	script, err := os.ReadFile(cfg.CqlDb.InitTableScriptPath)
	if err != nil {
		panic(fmt.Errorf("fatal error reading initialize cassandra cql table script: %v", err))
	}

	statements := strings.Split(string(script), ";")

	for _, statement := range statements {
		trimmedStatement := strings.TrimSpace(statement)
		if trimmedStatement != "" {
			err := s.Query(trimmedStatement).Exec()
			if err != nil {
				panic(fmt.Errorf("fatal error executing initialize cassandra cql table script statement: %v", err))
			}
		}
	}
}
