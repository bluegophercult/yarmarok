package service

import (
	"fmt"
	"time"
)

// Participant represents a participant of the application.
type Participant struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"createdAt"`
}

// ParticipantRequest is a request for creating a new/updated participant.
type ParticipantRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Note  string `json:"note"`
}

// ParticipantService is a service for participants.
type ParticipantService interface {
	Create(p *ParticipantRequest) (id string, err error)
	Edit(id string, p *ParticipantRequest) error
	Delete(id string) error
	List() ([]Participant, error)
}

// ParticipantStorage is a storage for participants.
//
//go:generate mockgen -destination=mock_participant_storage_test.go -package=service github.com/kaznasho/yarmarok/service ParticipantStorage
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
func (pm *ParticipantManager) Create(p *ParticipantRequest) (string, error) {
	if err := validateParticipant(p); err != nil {
		return "", err
	}

	prt := toParticipant(p)
	if err := pm.participantStorage.Create(prt); err != nil {
		return "", fmt.Errorf("creating participant: %w", err)
	}

	return prt.ID, nil
}

// Edit updates a participant.
func (pm *ParticipantManager) Edit(id string, p *ParticipantRequest) error {
	prt, err := pm.participantStorage.Get(id)
	if err != nil {
		return fmt.Errorf("getting participant: %w", err)
	}

	prt.Name = p.Name
	prt.Phone = p.Phone
	prt.Note = p.Note

	if err := pm.participantStorage.Update(prt); err != nil {
		return fmt.Errorf("updating participant: %w", err)
	}

	return nil
}

// Delete deletes a participant.
func (pm *ParticipantManager) Delete(id string) error {
	if err := pm.participantStorage.Delete(id); err != nil {
		return fmt.Errorf("deleting participant: %w", err)
	}

	return nil
}

// List returns all participants.
func (pm *ParticipantManager) List() ([]Participant, error) {
	prts, err := pm.participantStorage.GetAll()
	if err != nil {
		return nil, fmt.Errorf("getting all participants: %w", err)
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
