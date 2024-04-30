package database

import (
	"fmt"

	"spotigram/internal/config"
	infrastructureAbstractions "spotigram/internal/infrastructure/abstractions"

	"database/sql"

	"github.com/gocql/gocql"
	_ "github.com/lib/pq"
)

// Returns a postgres based sql database provider.
func NewCassandraCqlDatabaseProvider(cfg *config.Config) infrastructureAbstractions.CqlDatabaseProvider {
	cluster := gocql.NewCluster(cfg.CqlDb.Host)
	cluster.ProtoVersion = 4
	s, err := cluster.CreateSession()
	if err != nil {
		panic(fmt.Errorf("error creating cassandra session: %v", err))
	}
	initializeCassandraKeyspace(cfg, s)
	s.Close()

	cluster.Keyspace = cfg.CqlDb.KeySpace
	s, err = cluster.CreateSession()
	if err != nil {
		panic(fmt.Errorf("error creating cassandra session in keyspace: %v", err))
	}
	initializeCassandraTable(cfg, s)

	return &cassandraCqlDatabaseProvider{session: s}
}

// Returns a postgres based sql database provider.
func NewPostgresSqlDatabaseProvider(cfg *config.Config) infrastructureAbstractions.SqlDatabaseProvider {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.SqlDb.Host,
		cfg.SqlDb.User,
		cfg.SqlDb.Password,
		cfg.SqlDb.DBName,
		cfg.SqlDb.Port,
		cfg.SqlDb.SSLMode,
		cfg.SqlDb.TimeZone,
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

	return &postgresSqlDatabaseProvider{Db: db}
}
