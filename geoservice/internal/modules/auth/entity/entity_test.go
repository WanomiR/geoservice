package entity

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var auth *Auth

func init() {
	auth = NewAuth("localhost", "localhost", "secret", "localhost")
}

func TestAuth_ValidateEmail(t *testing.T) {
	testCases := []struct {
		name      string
		email     string
		wantError bool
	}{
		{"valid email", "email@example.com", false},
		{"valid email 2", "john.doe@gmail.com", false},
		{"invalid email", "email.com", true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if err := auth.ValidateEmail(testCase.email); (err != nil) != testCase.wantError {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, testCase.wantError)
			}
		})
	}
}

func TestAuth_ValidatePassword(t *testing.T) {
	testCases := []struct {
		name      string
		password  string
		wantError bool
	}{
		{"valid password", "password", false},
		{"valid password 2", ".m7?,,/", false},
		{"valid password 3", " ...77//~~~", false},
		{"invalid password", "em", true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if err := auth.ValidatePassword(testCase.password); (err != nil) != testCase.wantError {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, testCase.wantError)
			}
		})
	}
}

func TestAuth_VerifyPassword(t *testing.T) {
	t.Run("valid password", func(t *testing.T) {
		password := "password"
		encryptedPassword, _ := auth.EncryptPassword(password)

		if ok, err := auth.VerifyPassword(password, encryptedPassword); ok != true || err != nil {
			t.Errorf("VerifyPassword() error = %v, ok = %v, wantErr %v, wantOk %v", err, ok, false, true)
		}
	})

	t.Run("invalid password", func(t *testing.T) {
		password := "password"
		encryptedPassword, _ := auth.EncryptPassword(password)

		if ok, err := auth.VerifyPassword("otherPassword", encryptedPassword); !(ok != true || err != nil) {
			t.Errorf("VerifyPassword() error = %v, ok = %v, wantErr %v, wantOk %v", err, ok, false, false)
		}
	})
}

func TestAuth_VerifyRequest(t *testing.T) {
	// valid case with header
	token, _ := auth.GenerateToken("wanomir")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	wr := httptest.NewRecorder()
	req.Header.Set("Authorization", "Bearer "+token)

	t.Run("valid case with header", func(t *testing.T) {
		if _, _, err := auth.VerifyRequest(wr, req); err != nil {
			t.Errorf("VerifyRequest() error = %v, want %v", err, false)
		}
	})

	// invalid case with wrong header
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	wr = httptest.NewRecorder()
	req.Header.Set("Authorization", "Bearer")

	t.Run("invalid case with bad header", func(t *testing.T) {
		if _, _, err := auth.VerifyRequest(wr, req); err == nil {
			t.Errorf("VerifyRequest() error = %v, want %v", err, true)
		}
	})

	req.Header.Set("Authorization", "Bear "+token)
	t.Run("invalid case with bad header", func(t *testing.T) {
		if _, _, err := auth.VerifyRequest(wr, req); err == nil {
			t.Errorf("VerifyRequest() error = %v, want %v", err, true)
		}
	})

	// valid case with cookie
	cookie := auth.CreateCookie(token)
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	wr = httptest.NewRecorder()
	req.AddCookie(cookie)

	t.Run("valid case with cookie", func(t *testing.T) {
		if _, _, err := auth.VerifyRequest(wr, req); err != nil {
			t.Errorf("VerifyRequest() error = %v, want nil", err)
		}
	})

	// invalid case with expired cookie
	cookie = auth.CreateExpiredCookie()
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	wr = httptest.NewRecorder()
	req.AddCookie(cookie)

	t.Run("invalid case with expired cookie", func(t *testing.T) {
		if _, _, err := auth.VerifyRequest(wr, req); err == nil {
			t.Errorf("VerifyRequest() error = %v, want %v", err, true)
		}
	})
}

func TestAuth_RequireAuthorization(t *testing.T) {
	// create handler middleware
	requireAuth := auth.RequireAuthorization(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	// valid case with header
	token, _ := auth.GenerateToken("wanomir")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	// set authorization header
	req.Header.Set("Authorization", "Bearer "+token)

	wr := httptest.NewRecorder()
	requireAuth.ServeHTTP(wr, req)

	t.Run("valid case with header", func(t *testing.T) {
		if _, _, err := auth.VerifyRequest(wr, req); err != nil {
			t.Errorf("VerifyRequest() error = %v, want %v", err, false)
		}
	})

	// invalid case with wrong header
	token, _ = auth.GenerateToken("wanomir")
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	// set bad authorization header
	req.Header.Set("Authorization", "Bearer ")

	wr = httptest.NewRecorder()
	requireAuth.ServeHTTP(wr, req)

	t.Run("invalid case with bad header", func(t *testing.T) {
		if _, _, err := auth.VerifyRequest(wr, req); err == nil {
			t.Errorf("VerifyRequest() error = %v, want %v", err, true)
		}
	})
}
