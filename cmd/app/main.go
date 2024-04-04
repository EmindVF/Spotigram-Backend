package main

import (
	"spotigram/internal/cache"
	"spotigram/internal/config"
	"spotigram/internal/database"
	"spotigram/internal/infrastructure"
	infrastructureAbstractions "spotigram/internal/infrastructure/abstractions"
	"spotigram/internal/server"
	serverAbstractions "spotigram/internal/server/abstractions"
	serverConfig "spotigram/internal/server/config"
	serviceAbstractions "spotigram/internal/service/abstractions"
)

func main() {
	cfg := config.GetConfig()

	serverConfig.SetupConfig(&cfg)

	infrastructureAbstractions.DatabaseProviderInstance = database.NewPostgresDatabaseProvider(&cfg)
	defer infrastructureAbstractions.DatabaseProviderInstance.GetDb().Close()

	cache.ConnectRedis(&cfg)
	serverAbstractions.JWTCacheInstance = cache.NewJWTCache()

	serviceAbstractions.UserRepositoryInstance = infrastructure.NewSqlUserRepository()
	serviceAbstractions.ServerInstance = server.NewFiberServer()

	serviceAbstractions.ServerInstance.Start()
}
