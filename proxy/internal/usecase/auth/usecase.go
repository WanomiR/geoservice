package auth

import (
	"net/http"
)

type AuthProvider interface {
	Register(email, password, firstName, lastName, age string) (userId int, err error)
	Authorize(email, password string) (token string, cookie *http.Cookie, err error)
	ResetCookie() (cookie *http.Cookie, err error)
	VerifyToken(token string) (ok bool, err error)
}

type AuthUsecase struct {
	provider AuthProvider
}

func NewAuthUsecase(provider AuthProvider) *AuthUsecase {
	return &AuthUsecase{provider: provider}
}

func (a *AuthUsecase) Register(email, password, firstName, lastName, age string) (userId int, err error) {
	if userId, err = a.provider.Register(email, password, firstName, lastName, age); err != nil {
		return 0, err
	}
	return userId, nil
}

func (a *AuthUsecase) Authorize(email, password string) (token string, cookie *http.Cookie, err error) {
	if token, cookie, err = a.provider.Authorize(email, password); err != nil {
		return "", nil, err
	}
	return token, cookie, nil
}

func (a *AuthUsecase) ResetCookie() (cookie *http.Cookie, err error) {
	if cookie, err = a.provider.ResetCookie(); err != nil {
		return nil, err
	}
	return cookie, nil
}

func (a *AuthUsecase) VerifyToken(token string) (ok bool, err error) {
	if ok, err = a.provider.VerifyToken(token); err != nil {
		return false, err
	}
	return ok, nil
}
