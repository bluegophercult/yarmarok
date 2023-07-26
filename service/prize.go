package service

import (
	"errors"
	"time"
)

var (
	ErrPrizeAlreadyExists = errors.New("prize already exists")
	ErrPrizeNotFound      = errors.New("prize not found")
)

// Prize represents a prize of the application.
type Prize struct {
	ID          string
	Name        string
	TicketCost  int
	Description string
	CreatedAt   time.Time
}

// PrizeCreateRequest is a request for creating a new prize.
type PrizeCreateRequest struct {
	Name        string
	TicketCost  int
	Description string
}

// PrizeEditRequest is a request for updating a prize.
type PrizeEditRequest Prize

// PrizeListResult is a response for listing prizes.
type PrizeListResult struct {
	Prizes []Prize
}

// PrizeService is a service for prizes.
type PrizeService interface {
	Create(p *PrizeCreateRequest) (*CreateResult, error)
	Edit(p *PrizeEditRequest) (*Result, error)
	List() (*PrizeListResult, error)
}

// PrizeStorage is a storage for prizes.
type PrizeStorage interface {
	Create(*Prize) error
	Get(id string) (*Prize, error)
	Update(*Prize) error
	GetAll() ([]Prize, error)
	Delete(id string) error
}

// PrizeManager is an implementation of PrizeService.
type PrizeManager struct {
	prizeStorage PrizeStorage
}

// NewPrizeManager creates a new PrizeManager.
func NewPrizeManager(ps PrizeStorage) *PrizeManager {
	return &PrizeManager{prizeStorage: ps}
}

// Create creates a new prize
func (pm *PrizeManager) Create(p *PrizeCreateRequest) (*CreateResult, error) {
	prize := toPrize(p)
	if err := pm.prizeStorage.Create(prize); err != nil {
		return nil, err
	}

	return &CreateResult{ID: prize.ID}, nil
}

// Edit updates a Prize TODO: edit after donate representation
func (pm *PrizeManager) Edit(p *PrizeEditRequest) (*Result, error) {
	prize, err := pm.prizeStorage.Get(p.ID)
	if err != nil {
		return &Result{StatusError}, err
	}

	if err := pm.prizeStorage.Update(prize); err != nil {
		return &Result{StatusError}, err
	}

	return &Result{StatusSuccess}, nil
}

// List returns all prizes.
func (pm *PrizeManager) List() (*PrizeListResult, error) {
	prizes, err := pm.prizeStorage.GetAll()
	if err != nil {
		return nil, err
	}

	return &PrizeListResult{Prizes: prizes}, nil
}

func toPrize(p *PrizeCreateRequest) *Prize {
	return &Prize{
		ID:          stringUUID(),
		Name:        p.Name,
		TicketCost:  p.TicketCost,
		Description: p.Description,
		CreatedAt:   timeNow(),
	}
}
