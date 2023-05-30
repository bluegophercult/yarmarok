// Package service provides the business logic of the application.
package service

import (
	"errors"
)

var (
	// ErrUserAlreadyExists is returned when a user already exists.
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrYarmarokAlreadyExists is returned when a yarmarok already exists.
	ErrYarmarokAlreadyExists = errors.New("yarmarok already exists")
)

// User represents a user of the application.
type User struct {
	ID string
}

// UserStorage is a storage for users.
type UserStorage interface {
	Create(User) error
	Exists(id string) (bool, error)
	YarmarokStorage(userID string) YarmarokStorage
}

// UserService is a service for users.
type UserService interface {
	InitUserIfNotExists(id string) error
	YarmarokService(userID string) YarmarokService
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

// YarmarokService is a service for yarmaroks.
func (um *UserManager) YarmarokService(userID string) YarmarokService {
	return NewYarmarokManager(um.userStorage.YarmarokStorage(userID))
}
