package service

import (
	"errors"
	"time"
)

var (
	ErrDonationAlreadyExists = errors.New("donation already exists")
	ErrDonationNotFound      = errors.New("donation not found")
)

type Donation struct {
	ID            string    `json:"id"`
	PrizeID       string    `json:"prizeId"`
	ParticipantID string    `json:"participantId"`
	Amount        int       `json:"amount"`
	TicketNumber  int       `json:"ticketNumber"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"createdAt"`
}

type DonationAddRequest struct {
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

type DonationEditRequest Donation

type DonationListResult struct {
	Donations []Donation `json:"donations"`
}

type DonationService interface {
	AddDonation(d *DonationAddRequest) (*CreateResult, error)
	EditDonation(d *DonationEditRequest) (*Result, error)
	ListDonation() (*DonationListResult, error)
	// DeleteDonation????
}

type DonationStorage interface {
	Create(participantStorage ParticipantStorage, prizeStorage PrizeStorage, d *Donation) error
	Get(id string) (*Donation, error)
	Update(*Donation) error
	GetAll() ([]Donation, error)
	Delete(id string) error
}

type DonationManager struct {
	donationStorage    DonationStorage
	participantStorage ParticipantStorage
	prizeStorage       PrizeStorage
}

func NewDonationManager(ds DonationStorage, ps ParticipantStorage, pzs PrizeStorage) *DonationManager {
	return &DonationManager{
		donationStorage:    ds,
		participantStorage: ps,
		prizeStorage:       pzs,
	}
}

func (dm *DonationManager) AddDonation(d *DonationAddRequest) (*CreateResult, error) {
	donation := toDonation(d)
	if err := dm.donationStorage.Create(dm.participantStorage, dm.prizeStorage, donation); err != nil {
		return nil, err
	}

	return &CreateResult{ID: donation.ID}, nil
}

func (dm *DonationManager) EditDonation(d *DonationEditRequest) (*Result, error) {
	donation, err := dm.donationStorage.Get(d.ID)
	if err != nil {
		return &Result{StatusError}, err
	}

	if err := dm.donationStorage.Update(donation); err != nil {
		return &Result{StatusError}, err
	}

	return &Result{StatusSuccess}, nil
}

func (dm *DonationManager) ListDonation() (*DonationListResult, error) {
	donations, err := dm.donationStorage.GetAll()
	if err != nil {
		return nil, err
	}

	return &DonationListResult{Donations: donations}, nil
}

func toDonation(d *DonationAddRequest) *Donation {
	return &Donation{
		ID:          stringUUID(),
		Amount:      d.Amount,
		Description: d.Description,
		CreatedAt:   timeNow(),
	}
}
