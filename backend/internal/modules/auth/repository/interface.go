package repository

import (
	"proxy/internal/modules/auth/entities"
)

//go:generate mockgen -source=./interface.go -destination=../../../mocks/mock_auth_repo/mock_auth_repo.go
type DatabaseRepo interface {
	GetUserByEmail(string) (entities.User, error)
	InsertUser(entities.User) error
}
