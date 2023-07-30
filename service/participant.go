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

// IDResponse is a response by creating a new entity.
type IDResponse struct {
	ID string
}

// ParticipantRequest is a request for creating/updating a participant.
type ParticipantRequest struct {
	Name  string
	Phone string
	Note  string
}

// ParticipantListResponse is a response for listing participants.
type ParticipantListResponse struct {
	Participants []Participant
}

// ParticipantService is a service for participants.
//
//go:generate mockgen -destination=mocks/participant_service_mock.go -package=mocks github.com/kaznasho/yarmarok/service ParticipantService
type ParticipantService interface {
	Create(p *ParticipantRequest) (*CreateResult, error)
	Edit(id string, p *ParticipantRequest) error
	Delete(id string) error
	List() (*ParticipantListResponse, error)
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
func (pm *ParticipantManager) Create(p *ParticipantRequest) (*CreateResult, error) {
	participant := toParticipant(p)
	if err := pm.participantStorage.Create(participant); err != nil {
		return nil, err
	}

	return &CreateResult{participant.ID}, nil
}

// Edit updates a participant.
func (pm *ParticipantManager) Edit(id string, p *ParticipantRequest) error {
	participant, err := pm.participantStorage.Get(id)
	if err != nil {
		return err
	}

	participant.Name = p.Name
	participant.Phone = p.Phone
	participant.Note = p.Note

	return pm.participantStorage.Update(participant)
}

// Delete deletes a participant.
func (pm *ParticipantManager) Delete(id string) error {
	return pm.participantStorage.Delete(id)
}

// List returns all participants.
func (pm *ParticipantManager) List() (*ParticipantListResponse, error) {
	participants, err := pm.participantStorage.GetAll()
	if err != nil {
		return nil, err
	}

	return &ParticipantListResponse{Participants: participants}, nil
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
