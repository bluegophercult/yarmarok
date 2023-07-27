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

// ParticipantAddRequest is a request for creating a new participant.
type ParticipantAddRequest struct {
	Name  string
	Phone string
	Note  string
}

// ParticipantEditRequest is a request for updating a participant.
type ParticipantEditRequest Participant

// ParticipantListResult is a response for listing participants.
type ParticipantListResult struct {
	Participants []Participant
}

// ParticipantService is a service for participants.
type ParticipantService interface {
	Add(p *ParticipantAddRequest) (*InitResult, error)
	Edit(p *ParticipantEditRequest) (*Result, error)
	List() (*ParticipantListResult, error)
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

// Add creates a new participant
func (pm *ParticipantManager) Add(p *ParticipantAddRequest) (*InitResult, error) {
	participant := toParticipant(p)
	if err := pm.participantStorage.Create(participant); err != nil {
		return nil, err
	}

	return &InitResult{ID: participant.ID}, nil
}

// Edit updates a participant
func (pm *ParticipantManager) Edit(p *ParticipantEditRequest) (*Result, error) {
	participant, err := pm.participantStorage.Get(p.ID)
	if err != nil {
		return &Result{StatusError}, err
	}

	participant.Name = p.Name
	participant.Phone = p.Phone
	participant.Note = p.Note

	if err := pm.participantStorage.Update(participant); err != nil {
		return &Result{StatusError}, err
	}

	return &Result{StatusSuccess}, nil
}

// List returns all participants.
func (pm *ParticipantManager) List() (*ParticipantListResult, error) {
	participants, err := pm.participantStorage.GetAll()
	if err != nil {
		return nil, err
	}

	return &ParticipantListResult{Participants: participants}, nil
}

func toParticipant(p *ParticipantAddRequest) *Participant {
	return &Participant{
		ID:        stringUUID(),
		Name:      p.Name,
		Phone:     p.Phone,
		Note:      p.Note,
		CreatedAt: timeNow(),
	}
}
