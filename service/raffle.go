package service

import (
	"bytes"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	// StatusSuccess is a success status sent by the service to the client.
	StatusSuccess = "success"
	// StatusError is an error status sent by the service to the client.
	StatusError = "error"
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
	ID          string
	OrganizerID string
	Name        string
	Note        string
	CreatedAt   time.Time
}

// RaffleService is a service for raffles.
type RaffleService interface {
	Init(*RaffleInitRequest) (*InitResult, error)
	Get(id string) (*Raffle, error)
	List() (*RaffleListResponse, error)
	Export(id string) (*RaffleExportResponse, error)
	ParticipantService(id string) ParticipantService
	PrizeService(id string) PrizeService
}

// RaffleStorage is a storage for raffles.
type RaffleStorage interface {
	Create(*Raffle) error
	Get(id string) (*Raffle, error)
	GetAll() ([]Raffle, error)
	ParticipantStorage(id string) ParticipantStorage
	PrizeStorage(id string) PrizeStorage
}

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

// Init initializes a raffle.
func (rm *RaffleManager) Init(raf *RaffleInitRequest) (*InitResult, error) {
	raffle := Raffle{
		ID:        stringUUID(),
		Name:      raf.Name,
		Note:      raf.Note,
		CreatedAt: timeNow(),
	}

	err := rm.raffleStorage.Create(&raffle)
	if err != nil {
		return nil, err
	}

	return &InitResult{
		ID: raffle.ID,
	}, nil
}

// Get returns a raffle by id.
func (rm *RaffleManager) Get(id string) (*Raffle, error) {
	return rm.raffleStorage.Get(id)
}

// List lists raffles in organizer's scope.
func (rm *RaffleManager) List() (*RaffleListResponse, error) {
	raffles, err := rm.raffleStorage.GetAll()
	if err != nil {
		return nil, fmt.Errorf("get all raffles: %w", err)
	}

	return &RaffleListResponse{
		Raffles: raffles,
	}, nil
}

func (rm *RaffleManager) Export(id string) (*RaffleExportResponse, error) {
	raf, err := rm.Get(id)
	if err != nil {
		return nil, fmt.Errorf("get raffle: %w", err)
	}

	prtList, err := rm.ParticipantService(id).List()
	if err != nil {
		return nil, fmt.Errorf("get participants: %w", err)
	}

	przList, err := rm.PrizeService(id).List()
	if err != nil {
		return nil, fmt.Errorf("get prizes: %w", err)
	}

	excel := NewExcel()

	buf := new(bytes.Buffer)
	if err := excel.WriteExcel(buf, raf, prtList.Participants, przList.Prizes); err != nil {
		return nil, fmt.Errorf("write excel: %w", err)
	}

	resp := RaffleExportResponse{
		FileName: fmt.Sprintf("yarmarok_%s.xlsx", raf.ID),
		Data:     buf.Bytes(),
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

// RaffleInitRequest is a request for initializing a raffle.
type RaffleInitRequest struct {
	Name string
	Note string
}

// InitResult is a generic result of entity initialization.
type InitResult struct {
	ID string
}

// Result is a generic result with status.
type Result struct {
	Status string
}

// RaffleListResponse is a response for listing raffles.
type RaffleListResponse struct {
	Raffles []Raffle
}

// RaffleExportResponse is a response for exporting a raffle sub-collections.
type RaffleExportResponse struct {
	FileName string
	Data     []byte
}
