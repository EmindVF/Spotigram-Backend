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

	infrastructureAbstractions.SqlDatabaseProviderInstance =
		database.NewPostgresSqlDatabaseProvider(&cfg)
	defer infrastructureAbstractions.SqlDatabaseProviderInstance.GetDb().Close()

	infrastructureAbstractions.CqlDatabaseProviderInstance =
		database.NewCassandraCqlDatabaseProvider(&cfg)
	defer infrastructureAbstractions.CqlDatabaseProviderInstance.GetSession().Close()

	serviceAbstractions.UserRepositoryInstance =
		infrastructure.NewSqlUserRepository()

	serviceAbstractions.FriendRepositoryInstance =
		infrastructure.NewSqlFriendRepository()

	serviceAbstractions.FriendRequestRepositoryInstance =
		infrastructure.NewSqlFriendRequestRepository()

	serviceAbstractions.ChatRepositoryInstance =
		infrastructure.NewCqlChatRepository()

	cache.ConnectRedis(&cfg)
	serviceAbstractions.JWTCacheInstance = cache.NewJWTCache()
	defer cache.RedisClient.Close()

	serverConfig.SetupConfig(&cfg)
	serviceAbstractions.ServerInstance =
		server.NewFiberServer(cfg.App.RequestSizeLimit)

	serviceAbstractions.ServerInstance.Start()
}
