package infrastructure

import (
	infastructureAbstractions "spotigram/internal/infrastructure/abstractions"
	"spotigram/internal/infrastructure/repositories"
	serviceAbstractions "spotigram/internal/service/abstractions"
)

// Returns an sql user repository.
func NewSqlUserRepository() serviceAbstractions.UserRepository {
	return &repositories.SqlUserRepository{DBProvider: infastructureAbstractions.DatabaseProviderInstance}
}
