package service

import (
	"errors"
	"time"
)

var (
	ErrParticipantAlreadyExists = errors.New("participant already exists")
	ErrParticipantNotFound      = errors.New("participant not found")
)

// Participant represents a participant of the application.
type Participant struct {
	ID        string
	Name      string
	Phone     string
	Note      string
	CreatedAt time.Time
}

// ParticipantRequest is a request for creating a new participant.
type ParticipantRequest struct {
	Name  string
	Phone string
	Note  string
}

// ParticipantListResult is a response for listing participants.
type ParticipantListResult struct {
	Participants []Participant
}

// ParticipantService is a service for participants.
type ParticipantService interface {
	Create(*ParticipantRequest) (string, error)
	Edit(string, *ParticipantRequest) error
	Delete(string) error
	List() ([]Participant, error)
}

// ParticipantStorage is a storage for participants.
type ParticipantStorage interface {
	Create(*Participant) error
	Get(id string) (*Participant, error)
	Update(*Participant) error
	GetAll() ([]Participant, error)
	Delete(id string) error
}

// ParticipantManager is an implementation of ParticipantService.
type ParticipantManager struct {
	participantStorage ParticipantStorage
}

// NewParticipantManager creates a new ParticipantManager.
func NewParticipantManager(ps ParticipantStorage) *ParticipantManager {
	return &ParticipantManager{participantStorage: ps}
}

// Create creates a new participant.
func (pm *ParticipantManager) Create(new *ParticipantRequest) (string, error) {
	prt := toParticipant(new)
	if err := pm.participantStorage.Create(prt); err != nil {
		return "", err
	}

	return prt.ID, nil
}

// Edit updates a participant.
func (pm *ParticipantManager) Edit(id string, upd *ParticipantRequest) error {
	prt, err := pm.participantStorage.Get(id)
	if err != nil {
		return err
	}

	prt.Name = upd.Name
	prt.Phone = upd.Phone
	prt.Note = upd.Note

	if err := pm.participantStorage.Update(prt); err != nil {
		return err
	}

	return nil
}

// Delete deletes a participant.
func (pm *ParticipantManager) Delete(id string) error {
	return pm.participantStorage.Delete(id)
}

// List returns all participants.
func (pm *ParticipantManager) List() ([]Participant, error) {
	prts, err := pm.participantStorage.GetAll()
	if err != nil {
		return nil, err
	}

	return prts, nil
}

func toParticipant(p *ParticipantRequest) *Participant {
	return &Participant{
		ID:        stringUUID(),
		Name:      p.Name,
		Phone:     p.Phone,
		Note:      p.Note,
		CreatedAt: timeNow(),
	}
}
