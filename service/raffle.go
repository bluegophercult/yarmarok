package service

import (
	"bytes"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// stringUUID is a plumbing function for generating UUIDs.
// It is overridden in tests.
var stringUUID = func() string {
	return uuid.New().String()
}

// timeNow is a plumbing function for getting the current time.
// It is overridden in tests.
var timeNow = func() time.Time {
	return time.Now()
}

// Raffle represents a raffle.
type Raffle struct {
	ID          string    `json:"id"`
	OrganizerID string    `json:"organizerId"`
	Name        string    `json:"name"`
	Note        string    `json:"note"`
	CreatedAt   time.Time `json:"createdAt"`
}

// RaffleService is a service for raffles.
type RaffleService interface {
	Create(*RaffleRequest) (id string, err error)
	Get(id string) (*Raffle, error)
	Edit(id string, r *RaffleRequest) error
	Delete(id string) error
	List() ([]Raffle, error)
	Export(id string) (*RaffleExportResult, error)
	ParticipantService(id string) ParticipantService
	PrizeService(id string) PrizeService
}

// RaffleStorage is a storage for raffles.
//
//go:generate mockgen -destination=mock_raffle_storage_test.go -package=service github.com/kaznasho/yarmarok/service RaffleStorage
type RaffleStorage interface {
	Create(*Raffle) error
	Get(id string) (*Raffle, error)
	Update(*Raffle) error
	Delete(id string) error
	GetAll() ([]Raffle, error)
	ParticipantStorage(id string) ParticipantStorage
	PrizeStorage(id string) PrizeStorage
}

var _ RaffleService = (*RaffleManager)(nil)

// RaffleManager is an implementation of RaffleService.
type RaffleManager struct {
	raffleStorage RaffleStorage
}

// NewRaffleManager creates a new RaffleManager.
func NewRaffleManager(rs RaffleStorage) *RaffleManager {
	return &RaffleManager{
		raffleStorage: rs,
	}
}

// Create initializes a raffle.
func (rm *RaffleManager) Create(raf *RaffleRequest) (string, error) {
	if err := validateRaffle(raf); err != nil {
		return "", err
	}

	raffle := Raffle{
		ID:        stringUUID(),
		Name:      raf.Name,
		Note:      raf.Note,
		CreatedAt: timeNow(),
	}

	if err := rm.raffleStorage.Create(&raffle); err != nil {
		return "", err
	}

	return raffle.ID, nil
}

// Get returns a raffle by id.
func (rm *RaffleManager) Get(id string) (*Raffle, error) {
	return rm.raffleStorage.Get(id)
}

// Edit edits a raffle.
func (rm *RaffleManager) Edit(id string, r *RaffleRequest) error {
	raffle, err := rm.Get(id)
	if err != nil {
		return fmt.Errorf("get raffle: %w", err)
	}

	raffle.Name = r.Name
	raffle.Note = r.Note

	if err := rm.raffleStorage.Update(raffle); err != nil {
		return fmt.Errorf("update raffle: %w", err)
	}

	return nil
}

// Delete a raffle.
func (rm *RaffleManager) Delete(id string) error {
	if err := rm.raffleStorage.Delete(id); err != nil {
		return fmt.Errorf("deleting raffle: %w", err)
	}

	return nil
}

// List lists raffles in organizer's scope.
func (rm *RaffleManager) List() ([]Raffle, error) {
	raffles, err := rm.raffleStorage.GetAll()
	if err != nil {
		return nil, fmt.Errorf("get all raffles: %w", err)
	}

	return raffles, nil
}

func (rm *RaffleManager) Export(id string) (*RaffleExportResult, error) {
	raf, err := rm.Get(id)
	if err != nil {
		return nil, fmt.Errorf("get raffle: %w", err)
	}

	prts, err := rm.ParticipantService(id).List()
	if err != nil {
		return nil, fmt.Errorf("get participants: %w", err)
	}

	przs, err := rm.PrizeService(id).List()
	if err != nil {
		return nil, fmt.Errorf("get prizes: %w", err)
	}

	xlsx := NewXLSX()

	buf := new(bytes.Buffer)
	if err := xlsx.WriteXLSX(buf, raf, prts, przs); err != nil {
		return nil, fmt.Errorf("write xlsx: %w", err)
	}

	resp := RaffleExportResult{
		FileName: fmt.Sprintf("yarmarok_%s.xlsx", raf.ID),
		Content:  buf.Bytes(),
	}

	return &resp, nil
}

// ParticipantService is a service for participants.
func (rm *RaffleManager) ParticipantService(id string) ParticipantService {
	return NewParticipantManager(rm.raffleStorage.ParticipantStorage(id))
}

// PrizeService is a service for prizes.
func (rm *RaffleManager) PrizeService(id string) PrizeService {
	return NewPrizeManager(rm.raffleStorage.PrizeStorage(id))
}

// RaffleRequest is a request for initializing a raffle.
type RaffleRequest struct {
	Name string `json:"name" validate:"required,min=3,max=50,alphanumunicode,allowedChars"`
	Note string `json:"note" validate:"lte=1000"`
}

// RaffleExportResult is a response for exporting a raffle sub-collections.
type RaffleExportResult struct {
	FileName string `json:"fileName"`
	Content  []byte `json:"content"`
}
