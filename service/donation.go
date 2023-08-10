package service

import (
	"errors"
	"time"
)

// DonationService is a service for donations.
type DonationService interface {
	Create(*DonationRequest) (id string, err error)
	Get(id string) (*Donation, error)
	List() ([]Donation, error)
	Edit(id string, d *DonationRequest) error
	Delete(id string) error
}

// DonationStorage is a storage for donations.
//
//go:generate mockgen -destination=mock_donation_storage_test.go -package=service github.com/kaznasho/yarmarok/service DonationStorage
type DonationStorage interface {
	Create(*Donation) error
	Get(id string) (*Donation, error)
	GetAll() ([]Donation, error)
	Update(*Donation) error
	Delete(id string) error
}

// Donation represents a donation of the application.
type Donation struct {
	ID            string    `json:"id"`
	PrizeID       string    `json:"prizeId"`
	ParticipantID string    `json:"participantId"`
	Amount        int       `json:"amount"`
	TicketsNumber int       `json:"ticketsNumber"`
	CreatedAt     time.Time `json:"createdAt"`
}

// DonationRequest is a request for creating/updating a donation.
type DonationRequest struct {
	Amount        int    `json:"amount"`
	ParticipantID string `json:"participantId"`
}

var (
	ErrDonationAlreadyExists = errors.New("donation already exists")
	ErrDonationNotFound      = errors.New("donation not found")
)

var _ DonationService = (*DonationManager)(nil)

// DonationManager is an implementation of DonationService.
type DonationManager struct {
	donationStorage DonationStorage
	prizeStorage    PrizeStorage
}

// NewDonationManager creates a new DonationManager.
func NewDonationManager(ds DonationStorage, ps PrizeStorage) *DonationManager {
	return &DonationManager{
		donationStorage: ds,
		prizeStorage:    ps,
	}
}

// Create creates a new Donation.
func (dm *DonationManager) Create(d *DonationRequest) (string, error) {
	donation := toDonation(d)

	if err := dm.donationStorage.Create(donation); err != nil {
		return "", err
	}

	return donation.ID, nil
}

// Edit updates a Donation.
func (dm *DonationManager) Edit(id string, d *DonationRequest) error {
	donation, err := dm.donationStorage.Get(id)
	if err != nil {
		return err
	}

	donation.Amount = d.Amount
	donation.ParticipantID = d.ParticipantID

	if err := dm.donationStorage.Update(donation); err != nil {
		return err
	}

	return nil
}

// List returns a Donation list.
func (dm *DonationManager) List() ([]Donation, error) {
	donations, err := dm.donationStorage.GetAll()
	if err != nil {
		return nil, err
	}

	return donations, nil
}

// Get returns a Donation.
func (dm *DonationManager) Get(id string) (*Donation, error) {
	donation, err := dm.donationStorage.Get(id)
	if err != nil {
		return nil, err
	}

	return donation, nil
}

// Delete deletes a Donation.
func (dm *DonationManager) Delete(id string) error {
	if err := dm.donationStorage.Delete(id); err != nil {
		return err
	}

	return nil
}

func toDonation(d *DonationRequest) *Donation {
	return &Donation{
		ID:            stringUUID(),
		Amount:        d.Amount,
		ParticipantID: d.ParticipantID,
		CreatedAt:     timeNow(),
	}
}
