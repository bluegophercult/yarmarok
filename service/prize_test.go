package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type PrizeSuite struct {
	suite.Suite

	ctrl               *gomock.Controller
	storage            *MockPrizeStorage
	participantStorage *MockParticipantStorage

	manager *PrizeManager

	mockTime time.Time
	mockUUID string
}

func TestPrize(t *testing.T) {
	suite.Run(t, &PrizeSuite{})
}

func (s *PrizeSuite) SetupTest() {
	s.mockTime = time.Now().UTC()
	s.mockUUID = uuid.New().String()
	setTimeNowMock(s.mockTime)
	setUUIDMock(s.mockUUID)

	s.ctrl = gomock.NewController(s.T())
	s.storage = NewMockPrizeStorage(s.ctrl)
	s.participantStorage = NewMockParticipantStorage(s.ctrl)
	s.manager = NewPrizeManager(s.storage, s.participantStorage)
}

func (s *PrizeSuite) TestCreatePrize() {
	prizeRequest := dummyPrizeRequest()

	mockedPrize := &Prize{
		ID:          s.mockUUID,
		Name:        prizeRequest.Name,
		TicketCost:  prizeRequest.TicketCost,
		Description: prizeRequest.Description,
		CreatedAt:   s.mockTime,
	}

	s.storage.EXPECT().Create(mockedPrize).Return(nil)

	resID, err := s.manager.Create(prizeRequest)
	require.NoError(s.T(), err)
	require.Equal(s.T(), mockedPrize.ID, resID)

	s.Run("error", func() {
		s.storage.EXPECT().Create(mockedPrize).Return(assert.AnError)

		resID, err := s.manager.Create(prizeRequest)
		require.ErrorIs(s.T(), err, assert.AnError)
		require.Empty(s.T(), resID)
	})

	s.Run("invalid name", func() {
		request := dummyPrizeRequest()
		request.Name = "Ra"

		resID, err := s.manager.Create(request)
		require.Error(s.T(), err)
		require.Empty(s.T(), resID)
	})

	s.Run("invalid description", func() {
		request := dummyPrizeRequest()
		request.Description = "///"

		resID, err := s.manager.Create(request)
		require.Error(s.T(), err)
		require.Empty(s.T(), resID)
	})

	s.Run("invalid ticket cost", func() {
		request := dummyPrizeRequest()
		request.TicketCost = 0

		resID, err := s.manager.Create(request)
		require.Error(s.T(), err)
		require.Empty(s.T(), resID)
	})
}

func (s *PrizeSuite) TestGetPrize() {
	mockedPrize := dummyPrize()

	s.storage.EXPECT().Get(mockedPrize.ID).Return(mockedPrize, nil)

	res, err := s.manager.Get(mockedPrize.ID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), mockedPrize, res)

	s.Run("error", func() {
		s.storage.EXPECT().Get(mockedPrize.ID).Return(nil, assert.AnError)

		res, err := s.manager.Get(mockedPrize.ID)
		require.ErrorIs(s.T(), err, assert.AnError)
		require.Nil(s.T(), res)
	})
}

func (s *PrizeSuite) TestEditPrize() {
	prizeRequest := dummyPrizeRequest()
	mockedPrize := &Prize{
		ID:          s.mockUUID,
		Name:        prizeRequest.Name,
		TicketCost:  prizeRequest.TicketCost,
		Description: prizeRequest.Description,
		CreatedAt:   s.mockTime,
	}

	s.storage.EXPECT().Get(mockedPrize.ID).Return(mockedPrize, nil)
	s.storage.EXPECT().Update(mockedPrize).Return(nil)

	err := s.manager.Edit(mockedPrize.ID, prizeRequest)
	require.NoError(s.T(), err)

	s.Run("error", func() {
		s.storage.EXPECT().Get(mockedPrize.ID).Return(nil, assert.AnError)

		err := s.manager.Edit(mockedPrize.ID, prizeRequest)
		require.ErrorIs(s.T(), err, assert.AnError)
	})

	s.Run("error_in_update", func() {
		s.storage.EXPECT().Get(mockedPrize.ID).Return(mockedPrize, nil)
		s.storage.EXPECT().Update(mockedPrize).Return(assert.AnError)

		err := s.manager.Edit(mockedPrize.ID, prizeRequest)
		require.ErrorIs(s.T(), err, assert.AnError)
	})

	s.Run("invalid name", func() {
		request := dummyPrizeRequest()
		request.Name = "Ra"

		err := s.manager.Edit(mockedPrize.ID, request)
		require.Error(s.T(), err)
	})

	s.Run("invalid description", func() {
		request := dummyPrizeRequest()
		request.Description = "///"

		err := s.manager.Edit(mockedPrize.ID, request)
		require.Error(s.T(), err)
	})

	s.Run("invalid ticket cost", func() {
		request := dummyPrizeRequest()
		request.TicketCost = 0

		err := s.manager.Edit(mockedPrize.ID, request)
		require.Error(s.T(), err)
	})

	s.Run("already played", func() {
		prizeRequest := dummyPrizeRequest()
		mockedPrize := &Prize{
			ID:          s.mockUUID,
			Name:        prizeRequest.Name,
			TicketCost:  prizeRequest.TicketCost,
			Description: prizeRequest.Description,
			CreatedAt:   s.mockTime,
			PlayResult:  dummyPlayResult(),
		}

		s.storage.EXPECT().Get(mockedPrize.ID).Return(mockedPrize, nil)

		err := s.manager.Edit(mockedPrize.ID, prizeRequest)
		s.Require().ErrorIs(err, ErrPrizeAlreadyPlayed)
	})
}

func dummyPlayResult() *PrizePlayResult {
	return &PrizePlayResult{
		Winners: []PlayParticipant{
			*dummyplayParticipant(),
			*dummyplayParticipant(),
		},
		PlayParticipants: []PlayParticipant{
			*dummyplayParticipant(),
			*dummyplayParticipant(),
			*dummyplayParticipant(),
			*dummyplayParticipant(),
			*dummyplayParticipant(),
			*dummyplayParticipant(),
		},
	}
}

func dummyplayParticipant() *PlayParticipant {
	return &PlayParticipant{
		Participant: Participant{
			ID:        "participant_id_1",
			Name:      "participant_name_1",
			Phone:     "participant_phone_1",
			Note:      "participant_note_1",
			CreatedAt: time.Now().UTC(),
		},
		TotalDonation:      200,
		TotalTicketsNumber: 10,
		Donations: []Donation{
			{
				ID:            "donation_id_1",
				ParticipantID: "participant_id_1",
				Amount:        200,
				CreatedAt:     time.Now().UTC(),
			},
		},
	}
}

func (s *PrizeSuite) TestDeletePrize() {
	mockedPrize := dummyPrize()

	s.storage.EXPECT().Delete(mockedPrize.ID).Return(nil)

	err := s.manager.Delete(mockedPrize.ID)
	require.NoError(s.T(), err)

	s.Run("error", func() {
		s.storage.EXPECT().Delete(mockedPrize.ID).Return(assert.AnError)

		err := s.manager.Delete(mockedPrize.ID)
		require.ErrorIs(s.T(), err, assert.AnError)
	})
}

func (s *PrizeSuite) TestListPrize() {
	mockedPrize := dummyPrize()

	s.storage.EXPECT().GetAll().Return([]Prize{*mockedPrize}, nil)

	res, err := s.manager.List()
	require.NoError(s.T(), err)
	require.Equal(s.T(), []Prize{*mockedPrize}, res)

	s.Run("error", func() {
		s.storage.EXPECT().GetAll().Return(nil, assert.AnError)

		res, err := s.manager.List()
		require.ErrorIs(s.T(), err, assert.AnError)
		require.Nil(s.T(), res)
	})
}

func TestPlayPrize(t *testing.T) {
	suite.Run(t, &PlayPrizeSuite{})
}

type PlayPrizeSuite struct {
	PrizeSuite

	raffleID    string
	prizeID     string
	mockedPrize *Prize

	donationStorage *MockDonationStorage
}

func (s *PlayPrizeSuite) SetupTest() {
	s.PrizeSuite.SetupTest()

	s.raffleID = uuid.New().String()
	s.prizeID = uuid.New().String()
	s.mockedPrize = &Prize{
		ID:         s.prizeID,
		Name:       "Prize 1",
		TicketCost: 10,
	}

	s.donationStorage = NewMockDonationStorage(s.ctrl)

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

	s.participantStorage.EXPECT().GetAll().Return(participants, nil)
	s.storage.EXPECT().Get(s.prizeID).Return(s.mockedPrize, nil)
	s.donationStorage.EXPECT().GetAll().Return(donations, nil)
	// Pay attention, the mock always returns the first donation,
	// so winner in this test is always the same despite of the input.
	s.donationStorage.EXPECT().Get(MatcherAnyDonationID(donations...)).Return(&donations[0], nil)

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

	s.mockedPrize.PlayResult = &PrizePlayResult{
		Winners:          []PlayParticipant{expectedWinner},
		PlayParticipants: expectedParticipants,
	}

	s.storage.EXPECT().Update(s.mockedPrize).Return(nil)

	res, err := s.manager.Play(s.prizeID)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.Winners, 1)
	s.Require().Len(res.PlayParticipants, 2)

	s.Equal(expectedWinner, res.Winners[0])
	s.NotContains(res.PlayParticipants, expectedWinner)
}

func (s *PlayPrizeSuite) TestWinnerGeneration() {
	var r Randomizer = func(i uint) uint {
		return 0
	}

	s.Run("no_donations", func() {
		s.Panics(
			func() {
				donations := []Donation{}

				winnerID := r.GenerateWinner(donations, 100)
				s.Equal("", winnerID)
			},
		)
	})

	s.Run("zero_ticket_cost", func() {
		s.Panics(
			func() {
				donations := []Donation{
					{ID: "dn1", ParticipantID: "p1", Amount: 100},
				}
				ticketCost := 0

				r.GenerateWinner(donations, ticketCost)
			},
		)
	})

	s.Run("many_donations", func() {
		donations := []Donation{
			{ID: "dn1", ParticipantID: "p1", Amount: 100},
			{ID: "dn2", ParticipantID: "p1", Amount: 100},
			{ID: "dn3", ParticipantID: "p2", Amount: 200},
			{ID: "dn4", ParticipantID: "p2", Amount: 200},
			{ID: "dn5", ParticipantID: "p3", Amount: 300},
		}

		ticketCost := 100

		winnerID := r.GenerateWinner(donations, ticketCost)
		s.Equal("dn1", winnerID)
	})

	s.Run("one_donation", func() {
		donations := []Donation{
			{ID: "dn1", ParticipantID: "p1", Amount: 100},
			{ID: "dn2", ParticipantID: "p1", Amount: 100},
			{ID: "dn3", ParticipantID: "p2", Amount: 200},
			{ID: "dn4", ParticipantID: "p2", Amount: 200},
			{ID: "dn5", ParticipantID: "p3", Amount: 300},
		}

		ticketCost := 100

		winnerID := r.GenerateWinner(donations, ticketCost)
		s.Equal("dn1", winnerID)
	})
}

func (s *PlayPrizeSuite) TestGenerateDonationIDsList() {
	type testCase struct {
		donations  []Donation
		ticketCost int
		expected   []string
	}

	testCases := map[string]testCase{
		"no_donations": {
			donations:  []Donation{},
			ticketCost: 100,
			expected:   []string{},
		},
		"one_donation": {
			donations: []Donation{
				{ID: "dn1", ParticipantID: "p1", Amount: 100},
			},
			ticketCost: 100,
			expected:   []string{"dn1"},
		},
		"many_donations": {
			donations: []Donation{
				{ID: "dn1", ParticipantID: "p1", Amount: 100},
				{ID: "dn2", ParticipantID: "p1", Amount: 100},
				{ID: "dn3", ParticipantID: "p2", Amount: 200},
				{ID: "dn4", ParticipantID: "p2", Amount: 200},
				{ID: "dn5", ParticipantID: "p3", Amount: 300},
			},
			ticketCost: 100,
			expected: []string{
				"dn1",
				"dn2",
				"dn3",
				"dn3",
				"dn4",
				"dn4",
				"dn5",
				"dn5",
				"dn5",
			},
		},
		"separately_2_together_3": {
			donations: []Donation{
				{ID: "dn1", ParticipantID: "p1", Amount: 155},
				{ID: "dn1", ParticipantID: "p1", Amount: 155},
			},
			ticketCost: 100,
			expected: []string{
				"dn1",
				"dn1",
			},
		},
		"not_enough_money": {
			donations: []Donation{
				{ID: "dn1", ParticipantID: "p1", Amount: 50},
			},
			ticketCost: 100,
			expected:   []string{},
		},
		"almost_enough_money": {
			donations: []Donation{
				{ID: "dn1", ParticipantID: "p1", Amount: 99},
				{ID: "dn2", ParticipantID: "p1", Amount: 199},
			},
			ticketCost: 100,
			expected: []string{
				"dn2",
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			res := generateDonationIDsList(tc.donations, tc.ticketCost)
			s.Equal(tc.expected, res)
		})
	}

	s.Run("ticket_cost_0", func() {
		donations := []Donation{
			{ID: "dn1", ParticipantID: "p1", Amount: 100},
			{ID: "dn4", ParticipantID: "p2", Amount: 200},
			{ID: "dn5", ParticipantID: "p3", Amount: 300},
		}

		ticketCost := 0

		s.Panics(
			func() {
				generateDonationIDsList(donations, ticketCost)
			},
		)
	})
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

func dummyPrizeRequest() *PrizeRequest {
	return &PrizeRequest{
		Name:        "prize_name_1",
		TicketCost:  1234,
		Description: "prize_description_1",
	}
}

func dummyPrize() *Prize {
	return &Prize{
		ID:          "prize_id_1",
		Name:        "prize_name_1",
		TicketCost:  1234,
		Description: "prize_description_1",
		CreatedAt:   time.Now().UTC(),
	}
}
