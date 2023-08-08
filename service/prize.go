package service

import (
	"fmt"
	"time"
)

// Prize represents a prize of the application.
type Prize struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	TicketCost  int       `json:"ticketCost"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

// PrizeRequest is a request for creating a new prize.
type PrizeRequest struct {
	Name        string `json:"name"`
	TicketCost  int    `json:"ticketCost"`
	Description string `json:"description"`
}

// PrizeService is a service for prizes.
type PrizeService interface {
	Create(*PrizeRequest) (id string, err error)
	Get(id string) (*Prize, error)
	Edit(id string, p *PrizeRequest) error
	Delete(id string) error
	List() ([]Prize, error)
}

// PrizeStorage is a storage for prizes.

//go:generate mockgen -destination=mock_prize_storage_test.go -package=service github.com/kaznasho/yarmarok/service PrizeStorage
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
func (pm *PrizeManager) Create(p *PrizeRequest) (string, error) {
	prize := toPrize(p)
	if err := pm.prizeStorage.Create(prize); err != nil {
		return "", fmt.Errorf("create prize: %w", err)
	}

	return prize.ID, nil
}

// Get returns a Prize.
func (pm *PrizeManager) Get(id string) (*Prize, error) {
	prize, err := pm.prizeStorage.Get(id)
	if err != nil {
		return nil, fmt.Errorf("get prize: %w", err)
	}

	return prize, nil
}

// Edit updates a Prize.
func (pm *PrizeManager) Edit(id string, p *PrizeRequest) error {
	prize, err := pm.prizeStorage.Get(id)
	if err != nil {
		return fmt.Errorf("get prize: %w", err)
	}

	prize.Name = p.Name
	prize.TicketCost = p.TicketCost
	prize.Description = p.Description

	if err := pm.prizeStorage.Update(prize); err != nil {
		return fmt.Errorf("update prize: %w", err)
	}

	return nil
}

// Delete removes a Prize.
func (pm *PrizeManager) Delete(id string) error {
	if err := pm.prizeStorage.Delete(id); err != nil {
		return fmt.Errorf("delete prize: %w", err)
	}

	return nil
}

// List returns Prize list.
func (pm *PrizeManager) List() ([]Prize, error) {
	prizes, err := pm.prizeStorage.GetAll()
	if err != nil {
		return nil, fmt.Errorf("get all prizes: %w", err)
	}

	return prizes, nil
}

func toPrize(p *PrizeRequest) *Prize {
	return &Prize{
		ID:          stringUUID(),
		Name:        p.Name,
		TicketCost:  p.TicketCost,
		Description: p.Description,
		CreatedAt:   timeNow(),
	}
}
