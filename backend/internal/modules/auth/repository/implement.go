package repository

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"proxy/internal/modules/auth/entities"
	"sync"
)

var ErrorUserNotFound = errors.New("user not found")

type MapDBRepo struct {
	store map[string]string
	m     sync.RWMutex
}

func NewMapDBRepo(initUsers ...entities.User) *MapDBRepo {
	store := make(map[string]string)

	for _, user := range initUsers {
		password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		store[user.Email] = string(password)
	}

	return &MapDBRepo{store: store}
}

func (db *MapDBRepo) GetUserByEmail(userEmail string) (entities.User, error) {
	db.m.RLock() // blocks for writing
	defer db.m.RUnlock()

	for email, password := range db.store {
		if email == userEmail {
			return entities.User{Email: email, Password: password}, nil
		}
	}
	return entities.User{}, ErrorUserNotFound
}

func (db *MapDBRepo) InsertUser(user entities.User) error {
	db.m.Lock() // blocks for reading and writing
	defer db.m.Unlock()

	db.store[user.Email] = user.Password

	return nil
}
