// Package service provides the business logic of the application.
package service

import (
	"errors"
)

var (
	// ErrAlreadyExists is returned when a raffle already exists.
	ErrAlreadyExists = errors.New("item already exists")

	// ErrNotFound is returned when a raffle already exists.
	ErrNotFound = errors.New("item not found")
)

// Organizer represents an organizer of the application.
type Organizer struct {
	ID string `json:"id"`
}

// OrganizerStorage is a storage for organizers.
//
//go:generate mockgen -destination=mock_organizer_storage_test.go -package=service github.com/kaznasho/yarmarok/service OrganizerStorage
type OrganizerStorage interface {
	Create(*Organizer) error
	Exists(id string) (bool, error)
	RaffleStorage(organizerID string) RaffleStorage
}

// OrganizerService is a service for organizers.
type OrganizerService interface {
	CreateOrganizerIfNotExists(id string) error
	RaffleService(organizerID string) RaffleService
}

var _ OrganizerService = (*OrganizerManager)(nil)

// OrganizerManager is an implementation of OrganizerService.
type OrganizerManager struct {
	organizerStorage OrganizerStorage
}

// NewOrganizerManager creates a new OrganizerManager.
func NewOrganizerManager(os OrganizerStorage) *OrganizerManager {
	return &OrganizerManager{
		organizerStorage: os,
	}
}

// CreateOrganizerIfNotExists creates an organizer if it does not exist.
func (om *OrganizerManager) CreateOrganizerIfNotExists(id string) error {
	exists, err := om.organizerStorage.Exists(id)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	return om.organizerStorage.Create(&Organizer{ID: id})
}

// RaffleService is a service for raffles.
func (om *OrganizerManager) RaffleService(organizerID string) RaffleService {
	return NewRaffleManager(om.organizerStorage.RaffleStorage(organizerID))
}
