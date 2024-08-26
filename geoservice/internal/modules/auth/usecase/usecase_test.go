package usecase

import (
	"errors"
	"geoservice/internal/modules/auth/dto"
	"golang.org/x/crypto/bcrypt"
	"sync"
	"testing"
)

var authService *AuthService

func init() {
	authService = NewAuthService(
		"localhost", "localhost", "secret", "localhost",
		NewMockDBRepo(dto.User{"john.doe@gmail.com", "password"}),
	)
}

func TestAuthService_Register(t *testing.T) {
	type payload struct {
		email, password string
	}
	testCases := []struct {
		name      string
		payload   payload
		wantError bool
	}{
		{"successful registration", payload{"jen.star@gmail.com", "password"}, false},
		{"user already exists", payload{"john.doe@gmail.com", "password"}, true},
		{"invalid email", payload{"john.doe.gmail.com", "password"}, true},
		{"invalid password", payload{"john.doe@gmail.com", "pa"}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := authService.Register(tc.payload.email, tc.payload.password); (err != nil) != tc.wantError {
				t.Errorf("Register() error = %v, wantErr %v", err, tc.wantError)
			}
		})
	}
}

func TestAuthService_Authorize(t *testing.T) {
	type payload struct {
		email, password string
	}
	testCases := []struct {
		name      string
		payload   payload
		wantError bool
	}{
		{"successful authorization", payload{"john.doe@gmail.com", "password"}, false},
		{"user doesn't exist", payload{"new.user@gmail.com", "password"}, true},
		{"invalid password", payload{"john.doe@gmail.com", "pa"}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, err := authService.Authorize(tc.payload.email, tc.payload.password); (err != nil) != tc.wantError {
				t.Errorf("Authorize() error = %v, wantErr %v", err, tc.wantError)
			}
		})
	}
}

type MockDBRepo struct {
	store map[string]string
	m     *sync.RWMutex
}

func NewMockDBRepo(initUsers ...dto.User) *MockDBRepo {
	store := make(map[string]string)

	for _, user := range initUsers {
		password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		store[user.Email] = string(password)
	}

	return &MockDBRepo{
		store: store,
		m:     &sync.RWMutex{},
	}
}

func (db *MockDBRepo) GetUserByEmail(userEmail string) (dto.User, error) {
	db.m.RLock() // blocks for writing
	defer db.m.RUnlock()

	for email, password := range db.store {
		if email == userEmail {
			return dto.User{Email: email, Password: password}, nil
		}
	}
	return dto.User{}, errors.New("user not found")
}

func (db *MockDBRepo) InsertUser(email, password string) error {
	db.m.Lock() // blocks for reading and writing
	defer db.m.Unlock()

	db.store[email] = password

	return nil
}
