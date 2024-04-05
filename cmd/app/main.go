package main

import (
	"spotigram/internal/cache"
	"spotigram/internal/config"
	"spotigram/internal/database"
	"spotigram/internal/infrastructure"
	infrastructureAbstractions "spotigram/internal/infrastructure/abstractions"
	"spotigram/internal/server"
	serverConfig "spotigram/internal/server/config"
	serviceAbstractions "spotigram/internal/service/abstractions"
)

func main() {
	cfg := config.GetConfig()

	serverConfig.SetupConfig(&cfg)

	infrastructureAbstractions.DatabaseProviderInstance = database.NewPostgresDatabaseProvider(&cfg)
	defer infrastructureAbstractions.DatabaseProviderInstance.GetDb().Close()

	cache.ConnectRedis(&cfg)
	serviceAbstractions.JWTCacheInstance = cache.NewJWTCache()
	defer cache.RedisClient.Close()

	serviceAbstractions.UserRepositoryInstance = infrastructure.NewSqlUserRepository()
	serviceAbstractions.ServerInstance = server.NewFiberServer()

	serviceAbstractions.ServerInstance.Start()
}
