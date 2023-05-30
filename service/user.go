// Package service provides the business logic of the application.
package service

import (
	"errors"
)

var (
	// ErrUserAlreadyExists is returned when a user already exists.
	ErrUserAlreadyExists = errors.New("user already exists")
)

// User represents a user of the application.
type User struct {
	ID string
}

// UserStorage is a storage for users.
type UserStorage interface {
	Create(User) error
	Exists(id string) (bool, error)
}

// UserService is a service for users.
type UserService interface {
	InitUserIfNotExists(id string) error
}

// UserManager is an implementation of UserService.
type UserManager struct {
	userStorage UserStorage
}

// NewUserManager creates a new UserManager.
func NewUserManager(us UserStorage) *UserManager {
	return &UserManager{
		userStorage: us,
	}
}

// InitUserIfNotExists initializes a user if it does not exist.
func (um *UserManager) InitUserIfNotExists(id string) error {
	exists, err := um.userStorage.Exists(id)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	return um.userStorage.Create(User{ID: id})
}
