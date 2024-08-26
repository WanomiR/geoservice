package usecase

import (
	"errors"
	"geoservice/internal/modules/auth/dto"
	"geoservice/internal/modules/auth/entity"
	"github.com/wanomir/e"
	"net/http"
)

type User struct {
	Email, Password string
}

type DatabaseRepo interface {
	GetUserByEmail(email string) (dto.User, error)
	InsertUser(email, password string) error
}

type AuthService struct {
	auth *entity.Auth
	db   DatabaseRepo
}

func NewAuthService(issuer, audience, secret, cookieDomain string, dbRepo DatabaseRepo) *AuthService {
	return &AuthService{
		db:   dbRepo,
		auth: entity.NewAuth(issuer, audience, secret, cookieDomain),
	}
}

func (s *AuthService) Register(email, password string) error {
	if err := s.auth.ValidateEmail(email); err != nil {
		return err
	}

	if err := s.auth.ValidatePassword(password); err != nil {
		return err
	}

	if _, err := s.db.GetUserByEmail(email); err == nil {
		return errors.New("user already exists")
	}

	password, _ = s.auth.EncryptPassword(password)

	if err := s.db.InsertUser(email, password); err != nil {
		return e.Wrap("couldn't insert user", err)
	}

	return nil
}

func (s *AuthService) Authorize(email string, password string) (string, *http.Cookie, error) {
	user, err := s.db.GetUserByEmail(email)
	if err != nil {
		return "", nil, err
	}

	if ok, err := s.auth.VerifyPassword(password, user.Password); !ok || err != nil {
		return "", nil, errors.New("invalid password")
	}

	token, err := s.auth.GenerateToken(email)
	if err != nil {
		return "", nil, e.WrapIfErr("couldn't generate token", err)
	}

	cookie := s.auth.CreateCookie(token)

	return token, cookie, nil
}

func (s *AuthService) ResetCookie() *http.Cookie {
	return s.auth.CreateExpiredCookie()
}

func (s *AuthService) RequireAuthorization(next http.Handler) http.Handler {
	return s.auth.RequireAuthorization(next)
}
