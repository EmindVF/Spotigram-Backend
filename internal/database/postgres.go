package database

import (
	"fmt"
	"os"
	"strings"

	"spotigram/internal/config"

	"database/sql"

	_ "github.com/lib/pq"
)

type postgresDatabaseProvider struct {
	Db *sql.DB
}

func (p *postgresDatabaseProvider) GetDb() *sql.DB {
	return p.Db
}

func initializePostgresDatabase(cfg *config.Config, db *sql.DB) {
	script, err := os.ReadFile(cfg.Db.InitScriptPath)
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
