package http_v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/wanomir/rr"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"net/url"
	"proxy/internal/dto"
	"testing"
)

var cntrl *Controller

var mockUser = dto.User{
	Id:        1,
	Email:     "john.doe@gmail.com",
	Password:  "password",
	FirstName: "John",
	LastName:  "Doe",
	Age:       25,
}

func init() {
	cntrl = NewController(NewMockUseCase([]dto.User{mockUser}), rr.NewReadResponder(), &zap.Logger{})
}

func TestController_AddressSearch(t *testing.T) {
	testCases := []struct {
		name     string
		payload  RequestAddressSearch
		wantCode int
	}{
		{"normal request", RequestAddressSearch{"Подкопаевский переулок"}, 200},
		{"incomplete data", RequestAddressSearch{""}, 400},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var body bytes.Buffer
			json.NewEncoder(&body).Encode(tc.payload)

			req := httptest.NewRequest(http.MethodGet, "/address/search", &body)
			wr := httptest.NewRecorder()

			cntrl.AddressSearch(wr, req)

			r := wr.Result()

			if r.StatusCode != tc.wantCode {
				t.Errorf("got %d, want %d", r.StatusCode, tc.wantCode)
			}
		})
	}
}

func TestController_AddressGeocode(t *testing.T) {
	testCases := []struct {
		name     string
		payload  RequestAddressGeocode
		wantCode int
	}{
		{"normal request", RequestAddressGeocode{"55.753214", "37.642589"}, 200},
		{"incomplete data", RequestAddressGeocode{"", "37.642589"}, 400},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var body bytes.Buffer
			json.NewEncoder(&body).Encode(tc.payload)

			req := httptest.NewRequest(http.MethodGet, "/address/geocode", &body)
			wr := httptest.NewRecorder()

			cntrl.AddressGeocode(wr, req)

			r := wr.Result()

			if r.StatusCode != tc.wantCode {
				t.Errorf("got %d, want %d", r.StatusCode, tc.wantCode)
			}
		})
	}
}

func TestController_Register(t *testing.T) {
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
		{"existing user", payload{mockUser.Email, mockUser.Password}, 400},
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

			cntrl.Register(wr, req)

			r := wr.Result()

			if r.StatusCode != tc.wantStatus {
				t.Errorf("got status code %d, want %d", r.StatusCode, tc.wantStatus)
			}
		})
	}
}

func TestController_Login(t *testing.T) {
	type payload struct {
		email    string
		password string
	}
	testCases := []struct {
		name       string
		payload    payload
		wantStatus int
	}{
		{"successful login", payload{mockUser.Email, mockUser.Password}, 200},
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

			cntrl.Login(wr, req)

			r := wr.Result()

			if r.StatusCode != tc.wantStatus {
				t.Errorf("got status code %d, want %d", r.StatusCode, tc.wantStatus)
			}
		})
	}
}

func TestController_Logout(t *testing.T) {
	t.Run("successful logout", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/logout", nil)
		wr := httptest.NewRecorder()

		cntrl.Logout(wr, req)

		r := wr.Result()

		if r.StatusCode != 200 {
			t.Errorf("got status code %d, want 200", r.StatusCode)
		}
	})
}

func TestController_VerifyRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	wr := httptest.NewRecorder()
	req.Header.Set("Authorization", "Bearer token")

	t.Run("valid case with header", func(t *testing.T) {
		if err := cntrl.VerifyRequest(wr, req); err != nil {
			t.Errorf("VerifyRequest() error = %v, want %v", err, false)
		}
	})

	// invalid case with bad header
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	wr = httptest.NewRecorder()
	req.Header.Set("Authorization", "Bearer")

	t.Run("invalid case with bad header", func(t *testing.T) {
		if err := cntrl.VerifyRequest(wr, req); err == nil {
			t.Errorf("VerifyRequest() error = %v, want %v", err, true)
		}
	})

	req.Header.Set("Authorization", "Bear token")
	t.Run("invalid case with bad header", func(t *testing.T) {
		if err := cntrl.VerifyRequest(wr, req); err == nil {
			t.Errorf("VerifyRequest() error = %v, want %v", err, true)
		}
	})

}

type MockUseCase struct {
	mockUsers []dto.User
	token     string
}

func NewMockUseCase(users []dto.User) *MockUseCase {
	return &MockUseCase{
		mockUsers: users, token: "token",
	}
}

func (m *MockUseCase) AddressSearch(_ string) ([]dto.Address, error) {
	return []dto.Address{}, nil
}

func (m *MockUseCase) GeoCode(_, _ string) ([]dto.Address, error) {
	return []dto.Address{}, nil
}

func (m *MockUseCase) Register(email, password, firstName, lastName, age string) (int, error) {
	for _, u := range m.mockUsers {
		if email == u.Email {
			return 0, errors.New("user already exists")
		}
	}
	return len(m.mockUsers), nil
}

func (m *MockUseCase) Authorize(email string, password string) (string, *http.Cookie, error) {
	for _, u := range m.mockUsers {
		if email == u.Email && password == u.Password {
			return m.token, &http.Cookie{}, nil
		}
	}
	return "", nil, errors.New("invalid credentials")
}

func (m *MockUseCase) ResetCookie() (*http.Cookie, error) {
	return &http.Cookie{}, nil
}

func (m *MockUseCase) RequireAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (m *MockUseCase) VerifyToken(token string) (ok bool, err error) {
	if token == m.token {
		return true, nil
	}
	return false, errors.New("invalid token")
}
