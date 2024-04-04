package database

import (
	"fmt"

	"spotigram/internal/config"
	infrastructureAbstractions "spotigram/internal/infrastructure/abstractions"

	"database/sql"

	_ "github.com/lib/pq"
)

func NewPostgresDatabaseProvider(cfg *config.Config) infrastructureAbstractions.DatabaseProvider {
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

	return &postgresDatabaseProvider{Db: db}
}
