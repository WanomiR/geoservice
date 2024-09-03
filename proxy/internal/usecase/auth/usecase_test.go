package auth

import (
	"errors"
	"net/http"
	"proxy/internal/dto"
	"strconv"
	"testing"
)

var authUsecase *AuthUsecase

var mockUser = dto.User{
	Id:        1,
	Email:     "john.doe@gmail.com",
	Password:  "password",
	FirstName: "John",
	LastName:  "Doe",
	Age:       25,
}

func init() {
	authUsecase = NewAuthUsecase(NewMockProvider([]dto.User{mockUser}))
}

func TestAuthUsecase_Register(t *testing.T) {
	testCases := []struct {
		name    string
		payload dto.User
		wantErr bool
	}{
		{"successful registration", dto.User{Email: "jen.star@gmail.com", Password: "Password"}, false},
		{"existing user", mockUser, true},
	}

	for _, tc := range testCases {
		if _, err := authUsecase.Register(tc.payload.Email, tc.payload.Password, tc.payload.FirstName,
			tc.payload.LastName, strconv.Itoa(tc.payload.Age)); (err != nil) != tc.wantErr {
			t.Errorf("Register() error = %v, wantErr %v", err, tc.wantErr)
		}
	}
}

func TestAuthUsecase_Authorize(t *testing.T) {
	testCases := []struct {
		name            string
		email, password string
		wantErr         bool
	}{
		{"successful authorization", mockUser.Email, mockUser.Password, false},
		{"invalid credentials", "jen.star@gmail.com", "1234", true},
	}

	for _, tc := range testCases {
		if _, _, err := authUsecase.Authorize(tc.email, tc.password); (err != nil) != tc.wantErr {
			t.Errorf("Authorize() error = %v, wantErr %v", err, tc.wantErr)
		}
	}
}

func TestAuthUsecase_VerifyToken(t *testing.T) {
	t.Run("successful authorization", func(t *testing.T) {
		if ok, err := authUsecase.VerifyToken("token"); !ok || err != nil {
			t.Errorf("VerifyToken() ok = %v, error = %v, wantOk %v, wantErr %v", ok, err, true, false)
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		if ok, err := authUsecase.VerifyToken("1234"); ok || err == nil {
			t.Errorf("VerifyToken() ok = %v, error = %v, wantOk %v, wantErr %v", ok, err, false, true)
		}
	})
}

func TestAuthUsecase_ResetCookie(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		if _, err := authUsecase.ResetCookie(); err != nil {
			t.Errorf("ResetCookie() error = %v, wantErr %v", err, false)
		}
	})
}

type MockProvider struct {
	mockUsers []dto.User
	token     string
}

func NewMockProvider(users []dto.User) *MockProvider {
	return &MockProvider{mockUsers: users, token: "token"}
}

func (m *MockProvider) Register(email, _, _, _, _ string) (userId int, err error) {
	for _, user := range m.mockUsers {
		if user.Email == email {
			return 0, errors.New("user already exists")
		}
	}
	return len(m.mockUsers), nil
}

func (m *MockProvider) Authorize(email, password string) (token string, cookie *http.Cookie, err error) {
	for _, user := range m.mockUsers {
		if user.Email == email && user.Password == password {
			return m.token, new(http.Cookie), nil
		}
	}
	return "", nil, errors.New("invalid credentials")
}

func (m *MockProvider) ResetCookie() (cookie *http.Cookie, err error) {
	return new(http.Cookie), nil
}

func (m *MockProvider) VerifyToken(token string) (ok bool, err error) {
	if token == m.token {
		return true, nil
	}
	return false, errors.New("invalid token")
}
