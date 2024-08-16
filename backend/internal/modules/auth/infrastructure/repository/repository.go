package repository

import (
	"backend/internal/modules/auth/entity"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"sync"
)

type MapDBRepo struct {
	store map[string]string
	m     sync.RWMutex
}

func NewMapDBRepo(initUsers ...entity.User) *MapDBRepo {
	store := make(map[string]string)

	for _, user := range initUsers {
		password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		store[user.Email] = string(password)
	}

	return &MapDBRepo{store: store}
}

func (db *MapDBRepo) GetUserByEmail(userEmail string) (entity.User, error) {
	db.m.RLock() // blocks for writing
	defer db.m.RUnlock()

	for email, password := range db.store {
		if email == userEmail {
			return entity.User{Email: email, Password: password}, nil
		}
	}
	return entity.User{}, errors.New("user not found")
}

func (db *MapDBRepo) InsertUser(user entity.User) error {
	db.m.Lock() // blocks for reading and writing
	defer db.m.Unlock()

	db.store[user.Email] = user.Password

	return nil
}
