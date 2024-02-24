package service

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/exp/slices"
)

var (
	ErrPrizeAlreadyPlayed       = fmt.Errorf("prize already played")
	ErrNoParticipants           = fmt.Errorf("no participants")
	ErrNoDonations              = fmt.Errorf("no donations")
	ErrNotEnoughDonations       = fmt.Errorf("not enough donations")
	ErrEditPlayedPrizeDonations = fmt.Errorf("can't edit played prize donations")
)

// Prize represents a prize of the application.
type Prize struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	TicketCost  int              `json:"ticketCost"`
	Description string           `json:"description"`
	CreatedAt   time.Time        `json:"createdAt"`
	PlayResult  *PrizePlayResult `json:"playResult"`
}

// PrizePlayResult is a response for played prize
type PrizePlayResult struct {
	Winners          []PlayParticipant `json:"winners"`
	PlayParticipants []PlayParticipant `json:"participants"`
}

// PlayParticipant representation of result response of participant
type PlayParticipant struct {
	Participant        Participant `json:"participant"`
	TotalDonation      int         `json:"totalDonation"`
	TotalTicketsNumber int         `json:"totalTicketsNumber"`
	Donations          []Donation  `json:"donations"`
}

// PrizeRequest is a request for creating a new prize.
type PrizeRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=50,charsValidation"`
	TicketCost  int    `json:"ticketCost" validate:"gte=1,lte=5000"`
	Description string `json:"description" validate:"lte=1000,charsValidation"`
}

// Validate validates PrizeRequest.
func (p *PrizeRequest) Validate() error {
	return defaultValidator().Struct(p)
}

// PrizeService is a service for prizes.
type PrizeService interface {
	Create(*PrizeRequest) (id string, err error)
	Get(id string) (*Prize, error)
	Edit(id string, p *PrizeRequest) error
	Delete(id string) error
	List() ([]Prize, error)
	DonationService(id string) (DonationService, error)
	Play(prizeID string) (*PrizePlayResult, error)
}

// PrizeStorage is a storage for prizes.

//go:generate mockgen -destination=mock_prize_storage_test.go -package=service  github.com/bluegophercult/yarmarok/service PrizeStorage
type PrizeStorage interface {
	Create(*Prize) error
	Get(id string) (*Prize, error)
	Update(*Prize) error
	GetAll() ([]Prize, error)
	Delete(id string) error
	DonationStorage(id string) DonationStorage
}

// PrizeManager is an implementation of PrizeService.
type PrizeManager struct {
	prizeStorage       PrizeStorage
	participantStorage ParticipantStorage
	randomizer         Randomizer
}

// NewPrizeManager creates a new PrizeManager.
func NewPrizeManager(ps PrizeStorage, pts ParticipantStorage) *PrizeManager {
	return &PrizeManager{
		prizeStorage:       ps,
		participantStorage: pts,
		randomizer:         NewSimpleRandomizer(),
	}
}

// Create creates a new prize
func (pm *PrizeManager) Create(p *PrizeRequest) (string, error) {
	if err := p.Validate(); err != nil {
		return "", err
	}

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
	if err := p.Validate(); err != nil {
		return fmt.Errorf("validate prize: %w", err)
	}

	prize, err := pm.prizeStorage.Get(id)
	if err != nil {
		return fmt.Errorf("get prize: %w", err)
	}

	if prize.PlayResult != nil {
		return ErrPrizeAlreadyPlayed
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

// Play plays a prize.
func (pm *PrizeManager) Play(prizeID string) (*PrizePlayResult, error) {
	prize, err := pm.prizeStorage.Get(prizeID)
	if err != nil {
		return nil, fmt.Errorf("get prize to play: %w", err)
	}

	participants, err := pm.prepareParticipants(prize)
	if err != nil {
		return nil, fmt.Errorf("prepare participant for play: %w", err)
	}

	playResult := prize.Play(participants, pm.randomizer)
	err = pm.prizeStorage.Update(prize)
	if err != nil {
		return nil, fmt.Errorf("update prize with play results: %w", err)
	}

	return playResult, nil
}

func (pm *PrizeManager) prepareParticipants(prize *Prize) ([]PlayParticipant, error) {
	if prize.PlayResult != nil {
		if len(prize.PlayResult.PlayParticipants) == 0 {
			return nil, ErrNoParticipants
		}

		return prize.PlayResult.PlayParticipants, nil
	}

	participantList, err := pm.participantStorage.GetAll()
	if err != nil {
		return nil, fmt.Errorf("get participant list: %w", err)
	}

	if len(participantList) == 0 {
		return nil, ErrNoParticipants
	}

	ds := pm.prizeStorage.DonationStorage(prize.ID)
	donationsList, err := ds.GetAll()
	if err != nil {
		return nil, fmt.Errorf("get donation list: %w", err)
	}

	if len(donationsList) == 0 {
		return nil, ErrNoDonations
	}

	donations := countDonations(donationsList, participantList, prize.TicketCost)
	if len(donations) == 0 {
		return nil, ErrNotEnoughDonations
	}

	return donations, nil
}

func (p *Prize) Play(participants []PlayParticipant, randomizer Randomizer) *PrizePlayResult {
	winnerDonationID := randomizer.GenerateWinner(participants, p.TicketCost)

	winnerIndex := slices.IndexFunc(
		participants,
		func(p PlayParticipant) bool {
			return p.Participant.ID == winnerDonationID
		},
	)

	winner := participants[winnerIndex]
	participants = append(participants[:winnerIndex], participants[winnerIndex+1:]...)

	if p.PlayResult == nil {
		p.PlayResult = &PrizePlayResult{}
	}

	p.PlayResult.Winners = append(p.PlayResult.Winners, winner)
	p.PlayResult.PlayParticipants = participants

	return p.PlayResult
}

// countDonations counts donations, total amount and totat tickets count for each participant.
func countDonations(donations []Donation, participants []Participant, ticketCost int) []PlayParticipant {
	donationsMap := make(map[string][]Donation)

	for _, d := range donations {
		donationsMap[d.ParticipantID] = append(donationsMap[d.ParticipantID], d)
	}

	result := make([]PlayParticipant, 0, len(donationsMap))

	for _, participant := range participants {
		donations := donationsMap[participant.ID]
		totalDonation := countTotalDonation(donations)

		ticketsNumber := totalDonation / ticketCost
		if ticketsNumber == 0 {
			continue
		}

		result = append(result, PlayParticipant{
			Participant:        participant,
			TotalDonation:      totalDonation,
			TotalTicketsNumber: ticketsNumber,
			Donations:          donations,
		})
	}

	return result
}

func countTotalDonation(donations []Donation) int {
	total := 0

	for _, d := range donations {
		total += d.Amount
	}

	return total
}

// DonationService returns a DonationService for a prize.
func (pm *PrizeManager) DonationService(prizeID string) (DonationService, error) {
	prize, err := pm.prizeStorage.Get(prizeID)
	if err != nil {
		return nil, fmt.Errorf("get prize: %w", err)
	}

	donationStorage := pm.prizeStorage.DonationStorage(prize.ID)
	donationService := NewDonationManager(donationStorage)

	if prize.PlayResult != nil {
		return &ReadonlyDonationService{
			DonationService: donationService,
		}, nil
	}

	return donationService, nil
}

// ReadonlyDonationService is a DonationService that
// disallows editing donations for played prizes.
type ReadonlyDonationService struct {
	DonationService
}

func NewReadonlyDonationService(ds DonationService) *ReadonlyDonationService {
	return &ReadonlyDonationService{
		DonationService: ds,
	}
}

// Create is a stub that returns an error.
func (r *ReadonlyDonationService) Create(*DonationRequest) (string, error) {
	return "", ErrEditPlayedPrizeDonations
}

// Edit is a stub that returns an error.
func (r *ReadonlyDonationService) Edit(string, *DonationRequest) error {
	return ErrEditPlayedPrizeDonations
}

// Delete is a stub that returns an error.
func (r *ReadonlyDonationService) Delete(string) error {
	return ErrEditPlayedPrizeDonations
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

// Randomizer is a function type that returns a random number from 0 to n.
type Randomizer func(uint) uint

// GenerateWinner returns a winner ID.
// The winner is selected randomly as the person that made the donation.
// Function panics if donations list is empty or ticketCost is 0.
func (r Randomizer) GenerateWinner(participants []PlayParticipant, ticketCost int) (id string) {
	tickets := generateParticipantChanceList(participants, ticketCost)
	winnerTicketIndex := r(uint(len(tickets)))
	return tickets[winnerTicketIndex]
}

func generateParticipantChanceList(participants []PlayParticipant, ticketCost int) []string {
	participantIDs := make([]string, 0)
	for _, participant := range participants {
		for i := 0; i < participant.TotalTicketsNumber; i++ {
			participantIDs = append(participantIDs, participant.Participant.ID)
		}
	}

	return participantIDs
}

// NewSimpleRandomizer creates a new Randomizer
// that uses math/rand to generate random numbers.
func NewSimpleRandomizer() Randomizer {
	return func(i uint) uint {
		seed := time.Now().UnixNano()

		return uint(rand.New(rand.NewSource(seed)).Intn(int(i)))
	}
}
