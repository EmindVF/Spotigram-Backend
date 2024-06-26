package database

import (
	"fmt"
	"os"
	"strings"

	"spotigram/internal/config"

	"database/sql"

	_ "github.com/lib/pq"
)

type postgresSqlDatabaseProvider struct {
	Db *sql.DB
}

func (p *postgresSqlDatabaseProvider) GetDb() *sql.DB {
	return p.Db
}

// Initializes connected postgres database by the parameters in the config.
func initializePostgresDatabase(cfg *config.Config, db *sql.DB) {
	script, err := os.ReadFile(cfg.SqlDb.InitTableScriptPath)
	if err != nil {
		panic(fmt.Errorf("fatal error reading initialize postgres sql script: %v", err))
	}

	statements := strings.Split(string(script), ";")

	for _, statement := range statements {
		_, err := db.Exec(statement)
		if err != nil {
			panic(fmt.Errorf("fatal error executing initialize postgres sql script statement: %v", err))
		}
	}
}
