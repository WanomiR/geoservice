package repository

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"sync"
)

type User struct {
	Email, Password string
}

type MapDBRepo struct {
	store map[string]string
	m     *sync.RWMutex
}

func NewMapDBRepo(initUsers ...User) *MapDBRepo {
	store := make(map[string]string)

	for _, user := range initUsers {
		password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		store[user.Email] = string(password)
	}

	return &MapDBRepo{
		store: store,
		m:     &sync.RWMutex{},
	}
}

func (db *MapDBRepo) GetUserByEmail(userEmail string) (any, error) {
	db.m.RLock() // blocks for writing
	defer db.m.RUnlock()

	for email, password := range db.store {
		if email == userEmail {
			return User{Email: email, Password: password}, nil
		}
	}
	return User{}, errors.New("user not found")
}

func (db *MapDBRepo) InsertUser(email, password string) error {
	db.m.Lock() // blocks for reading and writing
	defer db.m.Unlock()

	db.store[email] = password

	return nil
}
