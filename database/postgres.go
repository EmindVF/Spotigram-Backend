package database

import (
	"fmt"
	"os"
	"strings"

	"spotigram/config"
	"spotigram/infrastructure"

	"database/sql"

	_ "github.com/lib/pq"
)

type postgresDatabase struct {
	Db *sql.DB
}

func (p *postgresDatabase) GetDb() *sql.DB {
	return p.Db
}

func NewPostgresDatabase(cfg *config.Config) infrastructure.DatabaseProvider {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.Db.Host,
		cfg.Db.User,
		cfg.Db.Password,
		cfg.Db.DBName,
		cfg.Db.Port,
		cfg.Db.SSLMode,
		cfg.Db.TimeZone,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic("failed to connect Postgres database.")
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("ping to database failed: %v", err))
	}

	initializePostgresDatabase(cfg, db)

	return &postgresDatabase{Db: db}
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
