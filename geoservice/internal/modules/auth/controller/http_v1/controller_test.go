package http_v1

import (
	"bytes"
	"errors"
	"geoservice/internal/modules/auth/entity"
	"github.com/wanomir/rr"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var MockUser = entity.User{
	Email:    "john.doe@gmail.com",
	Password: "password",
}

var auth *AuthController

func init() {
	mockAuthService := NewMockAuthService(MockUser)
	auth = NewAuthController(mockAuthService, rr.NewReadResponder())
}

func TestAuthController_Register(t *testing.T) {
	type payload struct {
		email    string
		password string
	}
	testCases := []struct {
		name       string
		payload    payload
		wantStatus int
	}{
		{"successful registry", payload{"some@user.com", "password"}, 201},
		{"existing user", payload{MockUser.Email, MockUser.Password}, 400},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			param := url.Values{}
			param.Set("email", tc.payload.email)
			param.Set("password", tc.payload.password)
			body := bytes.NewBufferString(param.Encode())

			req := httptest.NewRequest(http.MethodPost, "/api/register", body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			wr := httptest.NewRecorder()

			auth.Register(wr, req)

			r := wr.Result()

			if r.StatusCode != tc.wantStatus {
				t.Errorf("got status code %d, want %d", r.StatusCode, tc.wantStatus)
			}
		})
	}
}

func TestAuthController_Login(t *testing.T) {
	type payload struct {
		email    string
		password string
	}
	testCases := []struct {
		name       string
		payload    payload
		wantStatus int
	}{
		{"successful login", payload{MockUser.Email, MockUser.Password}, 200},
		{"unknown user", payload{"some@user.com", "password"}, 401},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			param := url.Values{}
			param.Set("email", tc.payload.email)
			param.Set("password", tc.payload.password)
			body := bytes.NewBufferString(param.Encode())

			req := httptest.NewRequest(http.MethodPost, "/api/login", body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			wr := httptest.NewRecorder()

			auth.Login(wr, req)

			r := wr.Result()

			if r.StatusCode != tc.wantStatus {
				t.Errorf("got status code %d, want %d", r.StatusCode, tc.wantStatus)
			}
		})
	}
}

type MockAuthService struct {
	mockUsers []entity.User
}

func TestAuthController_Logout(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/logout", nil)
	wr := httptest.NewRecorder()

	auth.Logout(wr, req)

	r := wr.Result()

	if r.StatusCode != 200 {
		t.Errorf("got status code %d, want %d", r.StatusCode, 200)
	}
}

func NewMockAuthService(mockUsers ...entity.User) *MockAuthService {
	a := new(MockAuthService)

	for _, user := range mockUsers {
		a.mockUsers = append(a.mockUsers, user)
	}

	return a
}

func (m *MockAuthService) Register(email, password string) error {
	for _, u := range m.mockUsers {
		if email == u.Email {
			return errors.New("user already exists")
		}
	}
	return nil
}

func (m *MockAuthService) Authorize(email string, password string) (string, *http.Cookie, error) {
	for _, u := range m.mockUsers {
		if email == u.Email && password == u.Password {
			return "token", &http.Cookie{}, nil
		}
	}
	return "", nil, errors.New("invalid credentials")
}

func (m *MockAuthService) ResetCookie() *http.Cookie {
	return &http.Cookie{}
}

func (m *MockAuthService) RequireAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
