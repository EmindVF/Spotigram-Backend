package infrastructure

import (
	infrastructureAbstractions "spotigram/internal/infrastructure/abstractions"
	"spotigram/internal/infrastructure/repositories"
	serviceAbstractions "spotigram/internal/service/abstractions"
)

// Returns an sql user repository.
func NewSqlUserRepository() serviceAbstractions.UserRepository {
	return &repositories.SqlUserRepository{
		DBProvider: infrastructureAbstractions.SqlDatabaseProviderInstance}
}

// Returns an sql friend repository.
func NewSqlFriendRepository() serviceAbstractions.FriendRepository {
	return &repositories.SqlFriendRepository{
		DBProvider: infrastructureAbstractions.SqlDatabaseProviderInstance}
}

// Returns an sql friend request repository.
func NewSqlFriendRequestRepository() serviceAbstractions.FriendRequestRepository {
	return &repositories.SqlFriendRequestRepository{
		DBProvider: infrastructureAbstractions.SqlDatabaseProviderInstance}
}

// Returns an cql chat repository.
func NewCqlChatRepository() serviceAbstractions.ChatRepository {
	return &repositories.CqlChatRepository{
		DBProvider: infrastructureAbstractions.CqlDatabaseProviderInstance}
}
