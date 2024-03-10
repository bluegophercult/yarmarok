package service

import (
	"fmt"
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
//go:generate mockgen -destination=mock_donation_storage_test.go -package=service  github.com/bluegophercult/yarmarok/service DonationStorage
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
	ParticipantID string    `json:"participantId"`
	Amount        int       `json:"amount"`
	CreatedAt     time.Time `json:"createdAt"`
}

// DonationRequest is a request for creating/updating a donation.
type DonationRequest struct {
	Amount        int    `json:"amount"`
	ParticipantID string `json:"participantId"`
}

var _ DonationService = (*DonationManager)(nil)

// DonationManager is an implementation of DonationService.
type DonationManager struct {
	donationStorage DonationStorage
}

// NewDonationManager creates a new DonationManager.
func NewDonationManager(ds DonationStorage) *DonationManager {
	return &DonationManager{
		donationStorage: ds,
	}
}

// Create creates a new Donation.
func (dm *DonationManager) Create(d *DonationRequest) (string, error) {
	donation := toDonation(d)
	return donation.ID, dm.donationStorage.Create(donation)
}

// Edit updates a Donation.
func (dm *DonationManager) Edit(id string, d *DonationRequest) error {
	donation, err := dm.donationStorage.Get(id)
	if err != nil {
		return fmt.Errorf("get donation: %w", err)
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
	return dm.donationStorage.GetAll()
}

// Get returns a Donation.
func (dm *DonationManager) Get(id string) (*Donation, error) {
	return dm.donationStorage.Get(id)
}

// Delete deletes a Donation.
func (dm *DonationManager) Delete(id string) error {
	return dm.donationStorage.Delete(id)
}

func toDonation(d *DonationRequest) *Donation {
	return &Donation{
		ID:            stringUUID(),
		Amount:        d.Amount,
		ParticipantID: d.ParticipantID,
		CreatedAt:     timeNow(),
	}
}
