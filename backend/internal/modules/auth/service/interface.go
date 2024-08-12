package service

import (
	"backend/internal/modules/auth/entities"
	"net/http"
)

//go:generate mockgen -source=./interface.go -destination=../../../mocks/mock_auth_service/mock_auth_service.go

type AuthServicer interface {
	EncryptPassword(password string) (string, error)
	VerifyPassword(password string, encryptedPassword string) (ok bool, err error)
	GenerateToken(subject string) (string, error)
	CreateCookie(refreshToken string) *http.Cookie
	CreateExpiredCookie() *http.Cookie
	VerifyRequest(w http.ResponseWriter, r *http.Request) (string, *entities.Claims, error)
}
