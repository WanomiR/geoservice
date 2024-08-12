package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var as *AuthService

func init() {
	as = NewAuthService("localhost", "localhost", "very-secret", "localhost")
}

func TestAuthService_VerifyPassword(t *testing.T) {
	password := "password"
	encrypted, _ := as.EncryptPassword(password)

	t.Run("normal case", func(t *testing.T) {
		if ok, _ := as.VerifyPassword(password, encrypted); !ok {
			t.Errorf("VerifyPassword() isValid = %v, want %v", false, true)
		}

		if _, err := as.VerifyPassword(encrypted, password); err == nil {
			t.Errorf("VerifyPassword() error = %v, want %v", err, true)
		}
	})
}

func TestAuthService_VerifyRequest(t *testing.T) {
	// valid case with header
	token, _ := as.GenerateToken("wanomir")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	wr := httptest.NewRecorder()
	req.Header.Set("Authorization", "Bearer "+token)

	t.Run("with token", func(t *testing.T) {
		if _, _, err := as.VerifyRequest(wr, req); err != nil {
			t.Errorf("VerifyRequest() error = %v, want nil", err)
		}
	})

	// invalid case with wrong header
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	wr = httptest.NewRecorder()
	req.Header.Set("Authorization", "Bearer")

	t.Run("invalid header", func(t *testing.T) {
		if _, _, err := as.VerifyRequest(wr, req); err == nil {
			t.Errorf("VerifyRequest() error = %v, want %v", err, true)
		}
	})

	req.Header.Set("Authorization", "Bear "+token)
	t.Run("invalid header", func(t *testing.T) {
		if _, _, err := as.VerifyRequest(wr, req); err == nil {
			t.Errorf("VerifyRequest() error = %v, want %v", err, true)
		}
	})

	// valid case with cookie
	cookie := as.CreateCookie(token)
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	wr = httptest.NewRecorder()
	req.AddCookie(cookie)

	t.Run("with cookie", func(t *testing.T) {
		if _, _, err := as.VerifyRequest(wr, req); err != nil {
			t.Errorf("VerifyRequest() error = %v, want nil", err)
		}
	})

	// invalid case with expired cookie
	cookie = as.CreateExpiredCookie()
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	wr = httptest.NewRecorder()
	req.AddCookie(cookie)

	t.Run("with expired cookie", func(t *testing.T) {
		if _, _, err := as.VerifyRequest(wr, req); err == nil {
			t.Errorf("VerifyRequest() error = %v, want %v", err, true)
		}
	})
}
