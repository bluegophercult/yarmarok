package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestPlayPrize(t *testing.T) {
	suite.Run(t, &PlayPrizeSuite{})
}

type PlayPrizeSuite struct {
	PrizeSuite

	raffleID string
	prizeID  string

	donationStorage *MockDonationStorage
}

func (s *PlayPrizeSuite) SetupTest() {
	s.PrizeSuite.SetupTest()

	s.raffleID = uuid.New().String()
	s.prizeID = uuid.New().String()

	s.donationStorage = NewMockDonationStorage(s.ctrl)
	s.manager.randomizer = func(i uint) uint {
		return 0
	}

	s.storage.EXPECT().DonationStorage(s.prizeID).Return(s.donationStorage).AnyTimes()
}

func (s *PlayPrizeSuite) TestPlayPrize() {
	participants := []Participant{
		{ID: "p1", Name: "Participant 1"},
		{ID: "p2", Name: "Participant 2"},
		{ID: "p3", Name: "Participant 3"},
	}

	donations := []Donation{
		{ID: "dn1", ParticipantID: "p1", Amount: 100},
		{ID: "dn2", ParticipantID: "p1", Amount: 100},
		{ID: "dn3", ParticipantID: "p2", Amount: 200},
		{ID: "dn4", ParticipantID: "p2", Amount: 200},
		{ID: "dn5", ParticipantID: "p3", Amount: 300},
	}

	mockedPrize := &Prize{
		ID:         s.prizeID,
		Name:       "Prize 1",
		TicketCost: 10,
	}

	s.participantStorage.EXPECT().GetAll().Return(participants, nil)
	s.storage.EXPECT().Get(s.prizeID).Return(mockedPrize, nil)
	s.donationStorage.EXPECT().GetAll().Return(donations, nil)

	expectedWinner := PlayParticipant{
		Participant:        participants[0],
		TotalDonation:      200,
		TotalTicketsNumber: 20,
		Donations:          donations[:2],
	}

	expectedParticipants := []PlayParticipant{
		{
			Participant:        participants[1],
			TotalDonation:      400,
			TotalTicketsNumber: 40,
			Donations:          donations[2:4],
		},
		{
			Participant:        participants[2],
			TotalDonation:      300,
			TotalTicketsNumber: 30,
			Donations:          donations[4:],
		},
	}

	expectedPrize := &Prize{
		ID:         s.prizeID,
		Name:       "Prize 1",
		TicketCost: 10,
		PlayResult: &PrizePlayResult{
			Winners:          []PlayParticipant{expectedWinner},
			PlayParticipants: expectedParticipants,
		},
	}

	s.storage.EXPECT().Update(expectedPrize).Return(nil)

	res, err := s.manager.Play(s.prizeID)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.Winners, 1)
	s.Require().Len(res.PlayParticipants, 2)

	s.Equal(expectedWinner, res.Winners[0])
	s.NotContains(res.PlayParticipants, expectedWinner)
}

func (s *PlayPrizeSuite) TestPlayPrizeNoParticipants() {
	s.storage.EXPECT().Get(s.prizeID).Return(dummyPrize(), nil)
	s.participantStorage.EXPECT().GetAll().Return([]Participant{}, nil)

	res, err := s.manager.Play(s.prizeID)
	s.Require().ErrorIs(err, ErrNoParticipants)
	s.Require().Nil(res)
}

func (s *PlayPrizeSuite) TestPlayPrizeNoDonations() {
	participants := dummyParticipantsList()

	mockPrize := dummyPrize()
	mockPrize.ID = s.prizeID

	s.participantStorage.EXPECT().GetAll().Return(participants, nil)
	s.storage.EXPECT().Get(s.prizeID).Return(mockPrize, nil)
	s.donationStorage.EXPECT().GetAll().Return([]Donation{}, nil)

	res, err := s.manager.Play(s.prizeID)
	s.Require().ErrorIs(err, ErrNoDonations)
	s.Require().Nil(res)
}

func (s *PlayPrizeSuite) TestPlayPrizeErrorListParticipants() {
	s.storage.EXPECT().Get(s.prizeID).Return(dummyPrize(), nil)
	s.participantStorage.EXPECT().GetAll().Return(nil, assert.AnError)

	res, err := s.manager.Play(s.prizeID)
	s.Require().ErrorIs(err, assert.AnError)
	s.Require().Nil(res)
}

func (s *PlayPrizeSuite) TestPlayPrizeErrorGetPrize() {
	s.storage.EXPECT().Get(s.prizeID).Return(nil, assert.AnError)

	res, err := s.manager.Play(s.prizeID)
	s.Require().ErrorIs(err, assert.AnError)
	s.Require().Nil(res)
}

func (s *PlayPrizeSuite) TestPlayPrizeErrorListDonations() {
	mockPrize := dummyPrize()
	mockPrize.ID = s.prizeID

	s.participantStorage.EXPECT().GetAll().Return(dummyParticipantsList(), nil)
	s.storage.EXPECT().Get(s.prizeID).Return(mockPrize, nil)
	s.donationStorage.EXPECT().GetAll().Return(nil, assert.AnError)

	res, err := s.manager.Play(s.prizeID)
	s.Require().ErrorIs(err, assert.AnError)
	s.Require().Nil(res)
}

func (s *PlayPrizeSuite) TestPlayPrizeErrorUpdatePrize() {
	participants := dummyParticipantsList()
	donations := []Donation{
		{ID: "dn1", ParticipantID: "p1", Amount: 100},
		{ID: "dn2", ParticipantID: "p1", Amount: 100},
		{ID: "dn3", ParticipantID: "p2", Amount: 200},
		{ID: "dn4", ParticipantID: "p2", Amount: 200},
		{ID: "dn5", ParticipantID: "p3", Amount: 300},
	}

	mockedPrize := &Prize{
		ID:         s.prizeID,
		Name:       "Prize 1",
		TicketCost: 10,
	}

	s.participantStorage.EXPECT().GetAll().Return(participants, nil)
	s.storage.EXPECT().Get(s.prizeID).Return(mockedPrize, nil)
	s.donationStorage.EXPECT().GetAll().Return(donations, nil)
	s.storage.EXPECT().Update(gomock.Any()).Return(assert.AnError)

	res, err := s.manager.Play(s.prizeID)
	s.Require().ErrorIs(err, assert.AnError)
	s.Require().Nil(res)
}

func (s *PlayPrizeSuite) TestPlayPrizeNotEnoughMoney() {
	participants := dummyParticipantsList()
	donations := []Donation{
		{ID: "dn1", ParticipantID: "p1", Amount: 50},
		{ID: "dn2", ParticipantID: "p1", Amount: 50},
		{ID: "dn3", ParticipantID: "p2", Amount: 200},
		{ID: "dn4", ParticipantID: "p2", Amount: 200},
		{ID: "dn5", ParticipantID: "p3", Amount: 300},
	}

	mockedPrize := &Prize{
		ID:         s.prizeID,
		Name:       "Prize 1",
		TicketCost: 1000,
	}

	s.participantStorage.EXPECT().GetAll().Return(participants, nil)
	s.storage.EXPECT().Get(s.prizeID).Return(mockedPrize, nil)
	s.donationStorage.EXPECT().GetAll().Return(donations, nil)

	res, err := s.manager.Play(s.prizeID)
	s.Require().ErrorIs(err, ErrNotEnoughDonations)
	s.Require().Nil(res)
}

func (s *PlayPrizeSuite) TestPlayPrizeSingleParticipant() {
	participants := dummyParticipantsList()
	donations := []Donation{
		{ID: "dn1", ParticipantID: "p1", Amount: 100},
		{ID: "dn2", ParticipantID: "p1", Amount: 100},
	}

	mockedPrize := &Prize{
		ID:         s.prizeID,
		Name:       "Prize 1",
		TicketCost: 10,
	}

	s.participantStorage.EXPECT().GetAll().Return(participants, nil)
	s.storage.EXPECT().Get(s.prizeID).Return(mockedPrize, nil)
	s.donationStorage.EXPECT().GetAll().Return(donations, nil)

	expectedWinner := PlayParticipant{
		Participant:        participants[0],
		TotalDonation:      200,
		TotalTicketsNumber: 20,
		Donations:          donations,
	}

	expectedPrize := &Prize{
		ID:         s.prizeID,
		Name:       "Prize 1",
		TicketCost: 10,
		PlayResult: &PrizePlayResult{
			Winners:          []PlayParticipant{expectedWinner},
			PlayParticipants: []PlayParticipant{},
		},
	}

	s.storage.EXPECT().Update(expectedPrize).Return(nil)

	res, err := s.manager.Play(s.prizeID)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.Winners, 1)
	s.Require().Len(res.PlayParticipants, 0)

	s.Equal(expectedWinner, res.Winners[0])
}

func (s *PlayPrizeSuite) TestPlayPrizeAgain() {
	participants := dummyParticipantsList()
	donations := []Donation{
		{ID: "dn1", ParticipantID: "p1", Amount: 100},
		{ID: "dn2", ParticipantID: "p1", Amount: 100},
		{ID: "dn3", ParticipantID: "p2", Amount: 200},
		{ID: "dn4", ParticipantID: "p2", Amount: 200},
		{ID: "dn5", ParticipantID: "p3", Amount: 300},
	}

	mockedPrize := &Prize{
		ID:         s.prizeID,
		Name:       "Prize 1",
		TicketCost: 10,
		PlayResult: &PrizePlayResult{
			Winners: []PlayParticipant{
				{
					Participant:        participants[0],
					TotalDonation:      200,
					TotalTicketsNumber: 20,
					Donations:          donations[:2],
				},
			},
			PlayParticipants: []PlayParticipant{
				{
					Participant:        participants[1],
					TotalDonation:      400,
					TotalTicketsNumber: 40,
					Donations:          donations[2:4],
				},
				{
					Participant:        participants[2],
					TotalDonation:      300,
					TotalTicketsNumber: 30,
					Donations:          donations[4:],
				},
			},
		},
	}

	expectedPrize := &Prize{
		ID:         s.prizeID,
		Name:       "Prize 1",
		TicketCost: 10,
		PlayResult: &PrizePlayResult{
			Winners: []PlayParticipant{
				{
					Participant:        participants[0],
					TotalDonation:      200,
					TotalTicketsNumber: 20,
					Donations:          donations[:2],
				},
				{
					Participant:        participants[1],
					TotalDonation:      400,
					TotalTicketsNumber: 40,
					Donations:          donations[2:4],
				},
			},
			PlayParticipants: []PlayParticipant{
				{
					Participant:        participants[2],
					TotalDonation:      300,
					TotalTicketsNumber: 30,
					Donations:          donations[4:],
				},
			},
		},
	}

	s.storage.EXPECT().Get(s.prizeID).Return(mockedPrize, nil)
	s.storage.EXPECT().Update(expectedPrize).Return(nil)

	res, err := s.manager.Play(s.prizeID)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	s.Equal(expectedPrize.PlayResult, res)
}

func (s *PlayPrizeSuite) TestPlayPrizeAgainNoParticipants() {
	mockedPrize := &Prize{
		ID:         s.prizeID,
		Name:       "Prize 1",
		TicketCost: 10,
		PlayResult: &PrizePlayResult{
			Winners: []PlayParticipant{
				{
					Participant:        Participant{ID: "p1"},
					TotalDonation:      200,
					TotalTicketsNumber: 20,
					Donations: []Donation{
						{ID: "dn1", ParticipantID: "p1", Amount: 100},
						{ID: "dn2", ParticipantID: "p1", Amount: 100},
					},
				},
			},
			PlayParticipants: []PlayParticipant{},
		},
	}

	s.storage.EXPECT().Get(s.prizeID).Return(mockedPrize, nil)

	res, err := s.manager.Play(s.prizeID)
	s.Require().ErrorIs(err, ErrNoParticipants)
	s.Require().Nil(res)
}

func (s *PlayPrizeSuite) TestPlayPrizeAgainNoDonations() {
	participants := dummyParticipantsList()

	mockedPrize := &Prize{
		ID:         s.prizeID,
		Name:       "Prize 1",
		TicketCost: 1000,
		PlayResult: &PrizePlayResult{
			Winners: []PlayParticipant{
				{
					Participant:        participants[0],
					TotalDonation:      200,
					TotalTicketsNumber: 20,
					Donations: []Donation{
						{ID: "dn1", ParticipantID: "p1", Amount: 1000},
						{ID: "dn2", ParticipantID: "p1", Amount: 1000},
					},
				},
			},
			PlayParticipants: []PlayParticipant{
				{
					Participant:        participants[1],
					TotalDonation:      400,
					TotalTicketsNumber: 0,
					Donations: []Donation{
						{ID: "dn3", ParticipantID: "p2", Amount: 200},
						{ID: "dn4", ParticipantID: "p2", Amount: 200},
					},
				},
				{
					Participant:        participants[2],
					TotalDonation:      300,
					TotalTicketsNumber: 0,
					Donations: []Donation{
						{ID: "dn5", ParticipantID: "p3", Amount: 300},
					},
				},
			},
		},
	}

	s.storage.EXPECT().Get(s.prizeID).Return(mockedPrize, nil)

	s.Panics(func() {
		// This scenario should simply never happen
		// because we can't have list of participants with not enough donations.
		s.manager.Play(s.prizeID)
	})
}

func (s *PlayPrizeSuite) TestPlayPrizeAgainErrorGetPrize() {
	s.storage.EXPECT().Get(s.prizeID).Return(nil, assert.AnError)

	res, err := s.manager.Play(s.prizeID)
	s.Require().ErrorIs(err, assert.AnError)
	s.Require().Nil(res)
}

func (s *PlayPrizeSuite) TestWinnerGeneration() {
	var r Randomizer = func(i uint) uint {
		return 0
	}

	s.Run("no_participants", func() {
		s.Panics(
			func() {
				participants := []PlayParticipant{}

				r.GenerateWinner(participants, 100)
			},
		)
	})

	s.Run("many_participants", func() {
		participants := []PlayParticipant{
			{
				Participant:        Participant{ID: "p1"},
				TotalDonation:      200,
				TotalTicketsNumber: 2,
				Donations: []Donation{
					{ID: "dn1", ParticipantID: "p1", Amount: 100},
					{ID: "dn2", ParticipantID: "p1", Amount: 100},
				},
			},
			{
				Participant:        Participant{ID: "p2"},
				TotalDonation:      400,
				TotalTicketsNumber: 4,
				Donations: []Donation{
					{ID: "dn3", ParticipantID: "p2", Amount: 200},
					{ID: "dn4", ParticipantID: "p2", Amount: 200},
				},
			},
			{
				Participant:        Participant{ID: "p3"},
				TotalDonation:      300,
				TotalTicketsNumber: 3,
				Donations: []Donation{
					{ID: "dn5", ParticipantID: "p3", Amount: 300},
				},
			},
		}

		ticketCost := 100

		winnerID := r.GenerateWinner(participants, ticketCost)
		s.Equal("p1", winnerID)
	})

	s.Run("one_participant", func() {
		participants := []PlayParticipant{
			{
				Participant:        Participant{ID: "p2"},
				TotalDonation:      400,
				TotalTicketsNumber: 4,
				Donations: []Donation{
					{ID: "dn3", ParticipantID: "p2", Amount: 200},
					{ID: "dn4", ParticipantID: "p2", Amount: 200},
				},
			},
		}

		ticketCost := 100

		winnerID := r.GenerateWinner(participants, ticketCost)
		s.Equal("p2", winnerID)
	})
}

func (s *PlayPrizeSuite) TestCountDonations() {
	type testCase struct {
		donations []Donation
		expected  []PlayParticipant
	}

	testCases := map[string]testCase{
		"no_donations": {
			donations: []Donation{},
			expected:  []PlayParticipant{},
		},
		"one_donation": {
			donations: []Donation{
				{ID: "dn1", ParticipantID: "p1", Amount: 100},
			},
			expected: []PlayParticipant{
				{
					Participant:        Participant{ID: "p1"},
					TotalDonation:      100,
					TotalTicketsNumber: 1,
					Donations: []Donation{
						{ID: "dn1", ParticipantID: "p1", Amount: 100},
					},
				},
			},
		},
		"many_donations": {
			donations: []Donation{
				{ID: "dn1", ParticipantID: "p1", Amount: 100},
				{ID: "dn2", ParticipantID: "p1", Amount: 100},
				{ID: "dn3", ParticipantID: "p2", Amount: 200},
				{ID: "dn4", ParticipantID: "p2", Amount: 200},
				{ID: "dn5", ParticipantID: "p3", Amount: 300},
			},
			expected: []PlayParticipant{
				{
					Participant:        Participant{ID: "p1"},
					TotalDonation:      200,
					TotalTicketsNumber: 2,
					Donations: []Donation{
						{ID: "dn1", ParticipantID: "p1", Amount: 100},
						{ID: "dn2", ParticipantID: "p1", Amount: 100},
					},
				},
				{
					Participant:        Participant{ID: "p2"},
					TotalDonation:      400,
					TotalTicketsNumber: 4,
					Donations: []Donation{
						{ID: "dn3", ParticipantID: "p2", Amount: 200},
						{ID: "dn4", ParticipantID: "p2", Amount: 200},
					},
				},
				{
					Participant:        Participant{ID: "p3"},
					TotalDonation:      300,
					TotalTicketsNumber: 3,
					Donations: []Donation{
						{ID: "dn5", ParticipantID: "p3", Amount: 300},
					},
				},
			},
		},
		"separately_2_together_3": {
			donations: []Donation{
				{ID: "dn1", ParticipantID: "p1", Amount: 155},
				{ID: "dn1", ParticipantID: "p1", Amount: 155},
			},
			expected: []PlayParticipant{
				{
					Participant:        Participant{ID: "p1"},
					TotalDonation:      310,
					TotalTicketsNumber: 3,
					Donations: []Donation{
						{ID: "dn1", ParticipantID: "p1", Amount: 155},
						{ID: "dn1", ParticipantID: "p1", Amount: 155},
					},
				},
			},
		},
		"not_enough_money": {
			donations: []Donation{
				{ID: "dn1", ParticipantID: "p1", Amount: 50},
			},
			expected: []PlayParticipant{},
		},
		"almost_enough_money": {
			donations: []Donation{
				{ID: "dn1", ParticipantID: "p1", Amount: 99},
				{ID: "dn2", ParticipantID: "p1", Amount: 199},
			},
			expected: []PlayParticipant{
				{
					Participant:        Participant{ID: "p1"},
					TotalDonation:      298,
					TotalTicketsNumber: 2,
					Donations: []Donation{
						{ID: "dn1", ParticipantID: "p1", Amount: 99},
						{ID: "dn2", ParticipantID: "p1", Amount: 199},
					},
				},
			},
		},
	}

	participants := []Participant{
		{ID: "p1"},
		{ID: "p2"},
		{ID: "p3"},
		{ID: "p4"},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			result := countDonations(tc.donations, participants, 100)
			s.Equal(tc.expected, result)
		})
	}
}

func MatcherAnyDonationID(donations ...Donation) gomock.Matcher {
	return gomock.Cond(func(donationID interface{}) bool {
		id := donationID.(string)
		for _, d := range donations {
			if d.ID == id {
				return true
			}
		}
		return false
	})
}

func dummyParticipantsList() []Participant {
	return []Participant{
		{ID: "p1", Name: "Participant 1"},
		{ID: "p2", Name: "Participant 2"},
		{ID: "p3", Name: "Participant 3"},
	}
}
