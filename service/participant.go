package service

import "errors"

var (
	ErrParticipantAlreadyExists = errors.New("participant already exists")
	ErrParticipantNotFound      = errors.New("participant not found")
)

// Participant represents a participant of the application.
type Participant struct {
	ID         string
	YarmarokID string
	Name       string
	Phone      string
	Email      string
	Notes      string
}

// ParticipantStorage is a storage for participants.
type ParticipantStorage interface {
	Create(Participant) error
	Get(id string) (*Participant, error)
	Update(Participant) error
	Delete(id string) error
	GetAll() ([]Participant, error)
	Exists(id string) (bool, error)
}

type ParticipantManager struct {
	participantStorage ParticipantStorage
}

func NewParticipantManager(ps ParticipantStorage) *ParticipantManager {
	return &ParticipantManager{participantStorage: ps}
}
