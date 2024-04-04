package infrastructure

import (
	infastructureAbstractions "spotigram/internal/infrastructure/abstractions"
	"spotigram/internal/infrastructure/repositories"
	serviceAbstractions "spotigram/internal/service/abstractions"
)

func NewSqlUserRepository() serviceAbstractions.UserRepository {
	return &repositories.SqlUserRepository{DBProvider: infastructureAbstractions.DatabaseProviderInstance}
}
