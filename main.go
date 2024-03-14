package main

import (
	"spotigram/config"
	"spotigram/database"
	"spotigram/infrastructure"
)

func main() {
	cfg := config.GetConfig()

	infrastructure.Database = database.NewPostgresDatabase(&cfg)

}
