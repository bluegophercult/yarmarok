package service

import (
	"time"
)

// Prize represents a prize of the application.
type Prize struct {
	ID          string
	Name        string
	TicketCost  int
	Description string
	CreatedAt   time.Time
}

// PrizeRequest is a request for creating a new prize.
type PrizeRequest struct {
	Name        string
	TicketCost  int
	Description string
}

// PrizeListResult is a response for listing prizes.
type PrizeListResult struct {
	Prizes []Prize
}

var _ PrizeService = (*PrizeManager)(nil)

// PrizeService is a service for prizes.
//
//go:generate mockgen -destination=../mocks/prize_service_mock.go -package=mocks github.com/kaznasho/yarmarok/service PrizeService
type PrizeService interface {
	Create(p *PrizeRequest) (*CreateResult, error)
	Edit(id string, p *PrizeRequest) error
	Delete(id string) error
	List() (*PrizeListResult, error)
}

// PrizeStorage is a storage for prizes.
//
//go:generate mockgen -destination=../mocks/prize_storage_mock.go -package=mocks github.com/kaznasho/yarmarok/service PrizeStorage
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
func (pm *PrizeManager) Create(p *PrizeRequest) (*CreateResult, error) {
	prize := toPrize(p)
	if err := pm.prizeStorage.Create(prize); err != nil {
		return nil, err
	}

	return &CreateResult{ID: prize.ID}, nil
}

// Edit updates a Prize.
func (pm *PrizeManager) Edit(id string, p *PrizeRequest) error {
	prize, err := pm.prizeStorage.Get(id)
	if err != nil {
		return err
	}

	prize.Name = p.Name
	prize.TicketCost = p.TicketCost
	prize.Description = p.Description

	return pm.prizeStorage.Update(prize)
}

func (pm *PrizeManager) Delete(id string) error {
	return pm.prizeStorage.Delete(id)
}

// List returns all prizes.
func (pm *PrizeManager) List() (*PrizeListResult, error) {
	prizes, err := pm.prizeStorage.GetAll()
	if err != nil {
		return nil, err
	}

	return &PrizeListResult{Prizes: prizes}, nil
}

func toPrize(p *PrizeRequest) *Prize {
	return &Prize{
		ID:          StringUUID(),
		Name:        p.Name,
		TicketCost:  p.TicketCost,
		Description: p.Description,
		CreatedAt:   TimeNow(),
	}
}
