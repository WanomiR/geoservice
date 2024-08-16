package usecase

import (
	"backend/internal/lib/e"
	"backend/internal/modules/auth/entity"
	"backend/internal/modules/auth/infrastructure/repository"
	"errors"
	"log"
	"net/http"
)

type DatabaseRepo interface {
	GetUserByEmail(email string) (entity.User, error)
	InsertUser(user entity.User) error
}

type AuthServicer interface {
	Register(user entity.User) error
	Authorize(email string, password string) (string, *http.Cookie, error)
	ResetCookie() *http.Cookie
	RequireAuthorization(next http.Handler) http.Handler
}

type AuthService struct {
	auth *entity.Auth
	db   DatabaseRepo
}

func NewAuthService(issuer, audience, secret, cookieDomain string) *AuthService {
	return &AuthService{
		db:   repository.NewMapDBRepo(entity.User{Email: "john.doe@gmail.com", Password: "password"}),
		auth: entity.NewAuth(issuer, audience, secret, cookieDomain),
	}
}

func (s *AuthService) Register(user entity.User) error {
	if len(user.Email) < 8 {
		return errors.New("email must be at least 7 characters")
	}

	if len(user.Password) < 8 {
		return errors.New("password must be at least 7 characters")
	}

	if _, err := s.db.GetUserByEmail(user.Email); err == nil {
		return errors.New("user already exists")
	}

	user.Password, _ = s.auth.EncryptPassword(user.Password)

	if err := s.db.InsertUser(user); err != nil {
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
		return "", nil, e.WrapIfErr("invalid password", err)
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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _, err := s.auth.VerifyRequest(w, r)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Authorization required"))

		} else {
			next.ServeHTTP(w, r)
		}
	})
}
