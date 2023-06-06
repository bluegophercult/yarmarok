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

// ParticipantService is a service for participants.
type ParticipantService interface {
	Add(p *ParticipantInitRequest) (*InitResult, error)
	Edit(p *ParticipantEditRequest) (*Response, error)
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

// Add creates a new participant
func (pm *ParticipantManager) Add(p *ParticipantInitRequest) (*InitResult, error) {
	participant := toParticipant(p)
	if err := pm.participantStorage.Create(participant); err != nil {
		return nil, err
	}

	return &InitResult{ID: participant.ID}, nil
}

// Edit updates a participant
func (pm *ParticipantManager) Edit(p *ParticipantEditRequest) (*Response, error) {
	participant, err := pm.participantStorage.Get(p.ID)
	if err != nil {
		return &Response{err.Error()}, err
	}

	if err := pm.participantStorage.Update(participant); err != nil {
		return &Response{err.Error()}, err
	}

	return &Response{"Successfully updated."}, nil
}

// List returns all participants.
func (pm *ParticipantManager) List() (*ParticipantListResponse, error) {
	participants, err := pm.participantStorage.GetAll()
	if err != nil {
		return nil, err
	}

	return &ParticipantListResponse{Participants: participants}, nil
}

// ParticipantInitRequest is a request for creating a new participant.
type ParticipantInitRequest struct {
	Name  string
	Phone string
	Note  string
}

// ParticipantEditRequest is a request for updating a participant.
type ParticipantEditRequest Participant

// ParticipantListResponse is a response for listing participants.
type ParticipantListResponse struct {
	Participants []Participant
}

// Response is a generic result.
type Response struct {
	Message string
}

func toParticipant(p *ParticipantInitRequest) *Participant {
	return &Participant{
		ID:        stringUUID(),
		Name:      p.Name,
		Phone:     p.Phone,
		Note:      p.Note,
		CreatedAt: timeNow(),
	}
}
