package service

import (
	"bytes"
	"fmt"
	"math/rand"
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
	PlayPrize(raffleID, prizeID string) (*PrizePlayResult, error)
	PlayPrizeAgain(raffleID, prizeID string, previousResult *PrizePlayResult) (*PrizePlayResult, error)
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

// PlayPrize find winner of prize
func (rm *RaffleManager) PlayPrize(raffleID, prizeID string) (*PrizePlayResult, error) {

	participantList, err := rm.ParticipantService(raffleID).List()
	if err != nil {
		return nil, fmt.Errorf("get participant list: %w", err)
	}

	pzs := rm.PrizeService(raffleID)
	prize, err := pzs.Get(prizeID)
	if err != nil {
		return nil, fmt.Errorf("get prize to play: %w", err)
	}

	ds := pzs.DonationService(prizeID)
	donationsList, err := ds.List()
	if err != nil {
		return nil, fmt.Errorf("get donation list: %w", err)
	}

	ticketCost := prize.TicketCost
	winnerDonationID := GetWinnerDonationID(donationsList, ticketCost)

	winnerDonation, err := ds.Get(winnerDonationID)
	if err != nil {
		return nil, fmt.Errorf("get winner donation: %w", err)
	}

	prizePlayResult := new(PrizePlayResult)
	for _, participant := range participantList {
		tempPlayParticipant := PlayParticipant{
			Participant: participant,
		}

		// add participant donations
		totalDonation := 0
		for _, donation := range donationsList {
			totalDonation += donation.Amount
			tempPlayParticipant.Donations = append(tempPlayParticipant.Donations, donation)
		}

		// calculate total donation and number of tickets
		tempPlayParticipant.TotalDonation = totalDonation
		tempPlayParticipant.NumberOfTickets = totalDonation / ticketCost

		if participant.ID == winnerDonation.ParticipantID {
			prizePlayResult.Winners = append(prizePlayResult.Winners, tempPlayParticipant)
			continue
		}

		prizePlayResult.PlayParticipants = append(prizePlayResult.PlayParticipants, tempPlayParticipant)
	}

	return prizePlayResult, nil
}

func (rm *RaffleManager) PlayPrizeAgain(raffleID, prizeID string, previousResult *PrizePlayResult) (*PrizePlayResult, error) {
	prize, err := rm.PrizeService(raffleID).Get(prizeID)
	if err != nil {
		return nil, fmt.Errorf("get prize to play: %w", err)
	}

	ticketCost := prize.TicketCost
	donations := make([]Donation, 0)
	for _, participant := range previousResult.PlayParticipants {
		donations = append(participant.Donations)
	}

	winnerDonationID := GetWinnerDonationID(donations, ticketCost)

	prizePlayResult := new(PrizePlayResult)

	// add previous winners
	for _, prewiousWinners := range previousResult.Winners {
		prizePlayResult.Winners = append(prizePlayResult.Winners, prewiousWinners)
	}

	skip := false
	for _, participant := range previousResult.PlayParticipants {
		// to skip adding winner to participants
		if skip == false {
			for _, donation := range participant.Donations {
				// add new winner
				if winnerDonationID == donation.ID {
					prizePlayResult.Winners = append(prizePlayResult.Winners, participant)
					skip = true
					break
				}
			}
			continue
		}

		prizePlayResult.PlayParticipants = append(prizePlayResult.PlayParticipants, participant)
	}

	return prizePlayResult, nil
}

// GetWinnerDonationID find donation that wins
func GetWinnerDonationID(donationsList []Donation, ticketCost int) (id string) {
	tickets := make([]string, 0)
	for _, donation := range donationsList {
		// calculate number of tickets in donation
		numberOfTickets := donation.Amount / ticketCost
		for i := 0; i < numberOfTickets; i++ {
			tickets = append(tickets, donation.ID)
		}
	}

	rand.Seed(time.Now().UnixNano())
	winnerDonationID := rand.Intn(len(tickets))

	return tickets[winnerDonationID]
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
	Name string `json:"name"`
	Note string `json:"note"`
}

// RaffleExportResult is a response for exporting a raffle sub-collections.
type RaffleExportResult struct {
	FileName string `json:"fileName"`
	Content  []byte `json:"content"`
}

// PrizePlayResult is a response for played prize
type PrizePlayResult struct {
	Winners          []PlayParticipant `json:"winners"`
	PlayParticipants []PlayParticipant `json:"participants"`
}

// PlayParticipant representation of result response of participant
type PlayParticipant struct {
	Participant     Participant `json:"participant"`
	TotalDonation   int         `json:"totalDonation"`
	NumberOfTickets int         `json:"numberOfTickets"`
	Donations       []Donation  `json:"donations"`
}
