// Package service provides the business logic of the application.
package service

import (
	"errors"
)

var (
	// ErrOrganizerAlreadyExists is returned when a organizer already exists.
	ErrOrganizerAlreadyExists = errors.New("organizer already exists")

	// ErrRaffleAlreadyExists is returned when a raffle already exists.
	ErrRaffleAlreadyExists = errors.New("raffle already exists")
)

// Organizer represents an organizer of the application.
type Organizer struct {
	ID string
}

// OrganizerStorage is a storage for organizers.
type OrganizerStorage interface {
	Create(Organizer) error
	Exists(id string) (bool, error)
	RaffleStorage(organizerID string) RaffleStorage
}

// OrganizerService is a service for organizers.
type OrganizerService interface {
	InitOrganizerIfNotExists(id string) error
	RaffleService(organizerID string) RaffleService
}

// OrganizerManager is an implementation of OrganizerService.
type OrganizerManager struct {
	organizerStorage OrganizerStorage
}

// NewOrganizerManager creates a new OrganizerManager.
func NewOrganizerManager(s OrganizerStorage) *OrganizerManager {
	return &OrganizerManager{
		organizerStorage: s,
	}
}

// InitOrganizerIfNotExists initializes an organizer if it does not exist.
func (m *OrganizerManager) InitOrganizerIfNotExists(id string) error {
	exists, err := m.organizerStorage.Exists(id)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	return m.organizerStorage.Create(Organizer{ID: id})
}

// RaffleService is a service for raffles.
func (m *OrganizerManager) RaffleService(organizerID string) RaffleService {
	return NewRaffleManager(m.organizerStorage.RaffleStorage(organizerID))
}
