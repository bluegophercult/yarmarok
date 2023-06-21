// Package service provides the business logic of the application.
package service

import (
	"errors"
)

var (
	// ErrOrganizerAlreadyExists is returned when a organizer already exists.
	ErrOrganizerAlreadyExists = errors.New("organizer already exists")

	// ErrYarmarokAlreadyExists is returned when a yarmarok already exists.
	ErrYarmarokAlreadyExists = errors.New("yarmarok already exists")
)

// Organizer represents a organizer of the application.
type Organizer struct {
	ID string
}

// OrganizerStorage is a storage for organizers.
type OrganizerStorage interface {
	Create(Organizer) error
	Exists(id string) (bool, error)
	YarmarokStorage(organizerID string) YarmarokStorage
}

// OrganizerService is a service for organizers.
type OrganizerService interface {
	InitOrganizerIfNotExists(id string) error
	YarmarokService(organizerID string) YarmarokService
}

// OrganizerManager is an implementation of OrganizerService.
type OrganizerManager struct {
	organizerStorage OrganizerStorage
}

// NewOrganizerManager creates a new OrganizerManager.
func NewOrganizerManager(us OrganizerStorage) *OrganizerManager {
	return &OrganizerManager{
		organizerStorage: us,
	}
}

// InitOrganizerIfNotExists initializes a organizer if it does not exist.
func (om *OrganizerManager) InitOrganizerIfNotExists(id string) error {
	exists, err := om.organizerStorage.Exists(id)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	return om.organizerStorage.Create(Organizer{ID: id})
}

// YarmarokService is a service for yarmaroks.
func (om *OrganizerManager) YarmarokService(organizerID string) YarmarokService {
	return NewYarmarokManager(om.organizerStorage.YarmarokStorage(organizerID))
}
