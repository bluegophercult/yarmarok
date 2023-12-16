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

type RaffleSuite struct {
	suite.Suite

	ctrl     *gomock.Controller
	storage  *MockRaffleStorage
	manager  *RaffleManager
	mockTime time.Time
	mockUUID string
}

func TestRaffle(t *testing.T) {
	suite.Run(t, &RaffleSuite{})
}

func (s *RaffleSuite) SetupTest() {
	s.mockTime = time.Now().UTC()
	s.mockUUID = uuid.New().String()
	setTimeNowMock(s.mockTime)
	setUUIDMock(s.mockUUID)

	s.ctrl = gomock.NewController(s.T())
	s.storage = NewMockRaffleStorage(s.ctrl)
	s.manager = NewRaffleManager(s.storage)
}

func (s *RaffleSuite) TestCreateRaffle() {
	raffleRequest := dummyRaffleRequest()

	mockedRaffle := Raffle{
		ID:        s.mockUUID,
		Name:      raffleRequest.Name,
		Note:      raffleRequest.Note,
		CreatedAt: s.mockTime,
	}

	s.storage.EXPECT().Create(&mockedRaffle).Return(nil)

	resID, err := s.manager.Create(raffleRequest)
	require.NoError(s.T(), err)
	require.Equal(s.T(), mockedRaffle.ID, resID)

	s.Run("error", func() {
		request := dummyRaffleRequest()
		expectedRaffle := &Raffle{
			ID:        s.mockUUID,
			Name:      request.Name,
			Note:      request.Note,
			CreatedAt: s.mockTime,
		}

		mockedErr := assert.AnError
		s.storage.EXPECT().Create(expectedRaffle).Return(mockedErr)

		response, err := s.manager.Create(request)
		s.ErrorIs(err, mockedErr)
		s.Equal("", response)
	})

	s.Run("invalid_name", func() {
		request := dummyRaffleRequest()
		request.Name = ""

		response, err := s.manager.Create(request)
		s.ErrorIs(err, ErrInvalidRequest)
		s.Equal("", response)
	})

	s.Run("invalid_note", func() {
		request := dummyRaffleRequest()
		request.Note = "<>////"

		response, err := s.manager.Create(request)
		s.ErrorIs(err, ErrInvalidRequest)
		s.Equal("", response)
	})
}

func (s *RaffleSuite) TestGetRaffle() {
	mockedRaffle := dummyRaffle()

	s.storage.EXPECT().Get(mockedRaffle.ID).Return(mockedRaffle, nil)

	res, err := s.manager.Get(mockedRaffle.ID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), mockedRaffle, res)

	s.Run("error", func() {
		mockedErr := assert.AnError
		s.storage.EXPECT().Get(mockedRaffle.ID).Return(nil, mockedErr)

		res, err := s.manager.Get(mockedRaffle.ID)
		s.ErrorIs(err, mockedErr)
		s.Nil(res)
	})
}

func (s *RaffleSuite) TestEditRaffle() {
	mockedRaffle := dummyRaffle()

	s.storage.EXPECT().Get(mockedRaffle.ID).Return(mockedRaffle, nil)
	s.storage.EXPECT().Update(mockedRaffle).Return(nil)

	err := s.manager.Edit(mockedRaffle.ID, dummyRaffleRequest())
	require.NoError(s.T(), err)

	s.Run("error", func() {
		mockedErr := assert.AnError
		s.storage.EXPECT().Get(mockedRaffle.ID).Return(nil, mockedErr)

		err := s.manager.Edit(mockedRaffle.ID, dummyRaffleRequest())
		s.ErrorIs(err, mockedErr)
	})

	s.Run("error_in_update", func() {
		mockedErr := assert.AnError
		s.storage.EXPECT().Get(mockedRaffle.ID).Return(mockedRaffle, nil)
		s.storage.EXPECT().Update(mockedRaffle).Return(mockedErr)

		err := s.manager.Edit(mockedRaffle.ID, dummyRaffleRequest())
		s.ErrorIs(err, mockedErr)
	})

	s.Run("invalid_name", func() {
		request := dummyRaffleRequest()
		request.Name = ""

		err := s.manager.Edit(mockedRaffle.ID, request)
		s.ErrorIs(err, ErrInvalidRequest)
	})

	s.Run("invalid_note", func() {
		request := dummyRaffleRequest()
		request.Note = "<>////"

		err := s.manager.Edit(mockedRaffle.ID, request)
		s.ErrorIs(err, ErrInvalidRequest)
	})
}

func (s *RaffleSuite) TestDeleteRaffle() {
	mockedRaffle := dummyRaffle()

	s.storage.EXPECT().Delete(mockedRaffle.ID).Return(nil)

	err := s.manager.Delete(mockedRaffle.ID)
	require.NoError(s.T(), err)

	s.Run("error", func() {
		mockedErr := assert.AnError
		s.storage.EXPECT().Delete(mockedRaffle.ID).Return(mockedErr)

		err := s.manager.Delete(mockedRaffle.ID)
		s.ErrorIs(err, mockedErr)
	})
}

func (s *RaffleSuite) TestListRaffles() {
	mockedRaffles := []Raffle{*dummyRaffle(), *dummyRaffle()}

	s.storage.EXPECT().GetAll().Return(mockedRaffles, nil)

	res, err := s.manager.List()
	require.NoError(s.T(), err)
	require.Equal(s.T(), mockedRaffles, res)

	s.Run("error", func() {
		mockedErr := assert.AnError
		s.storage.EXPECT().GetAll().Return(nil, mockedErr)

		res, err := s.manager.List()
		s.ErrorIs(err, mockedErr)
		s.Nil(res)
	})
}

func (s *RaffleSuite) TestExportRaffle() {
	raffle := &Raffle{ID: s.mockUUID, Name: "Raffle Test"}
	prts := []Participant{
		{ID: "p1", Name: "Participant 1"},
		{ID: "p2", Name: "Participant 2"},
	}
	przs := []Prize{
		{ID: "pr1", Name: "Prize 1"},
		{ID: "pr2", Name: "Prize 2"},
	}

	s.storage.EXPECT().Get(s.mockUUID).Return(raffle, nil)

	psMock := NewMockParticipantStorage(s.ctrl)
	s.storage.EXPECT().ParticipantStorage(s.mockUUID).Return(psMock)
	psMock.EXPECT().GetAll().Return(prts, nil)

	pzMock := NewMockPrizeStorage(s.ctrl)
	s.storage.EXPECT().PrizeStorage(s.mockUUID).Return(pzMock)
	pzMock.EXPECT().GetAll().Return(przs, nil)

	res, err := s.manager.Export(s.mockUUID)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal("yarmarok_"+s.mockUUID+".xlsx", res.FileName)
	s.Require().NotEmpty(res.Content)
}

func TestPlayPrize(t *testing.T) {
	suite.Run(t, &PlayPrizeSuite{})
}

type PlayPrizeSuite struct {
	RaffleSuite

	raffleID string
	prizeID  string

	mockedPrize *Prize

	participantStorage *MockParticipantStorage
	prizeStorage       *MockPrizeStorage
	donationStorage    *MockDonationStorage
}

func (s *PlayPrizeSuite) SetupTest() {
	s.RaffleSuite.SetupTest()

	s.raffleID = uuid.New().String()
	s.prizeID = uuid.New().String()
	s.mockedPrize = &Prize{
		ID:         s.prizeID,
		Name:       "Prize 1",
		TicketCost: 10,
	}

	s.participantStorage = NewMockParticipantStorage(s.ctrl)
	s.prizeStorage = NewMockPrizeStorage(s.ctrl)
	s.donationStorage = NewMockDonationStorage(s.ctrl)

	s.storage.EXPECT().ParticipantStorage(s.raffleID).Return(s.participantStorage).AnyTimes()
	s.storage.EXPECT().PrizeStorage(s.raffleID).Return(s.prizeStorage).AnyTimes()
	s.prizeStorage.EXPECT().DonationStorage(s.prizeID).Return(s.donationStorage).AnyTimes()
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
	s.prizeStorage.EXPECT().Get(s.prizeID).Return(s.mockedPrize, nil)
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

	s.prizeStorage.EXPECT().Update(s.mockedPrize).Return(nil)

	res, err := s.manager.PlayPrize(s.raffleID, s.prizeID)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.Winners, 1)
	s.Require().Len(res.PlayParticipants, 2)

	s.Equal(expectedWinner, res.Winners[0])
	s.NotContains(res.PlayParticipants, expectedWinner)
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

// func (s *PlayPrizeSuite) TestPlayPrizeAgain() {
// 	przs := []Prize{
// 		{ID: "pr1", Name: "Prize 1", TicketCost: 10},
// 		{ID: "pr2", Name: "Prize 2", TicketCost: 20},
// 	}

// 	mockedPreviousResult := &PrizePlayResult{
// 		Winners: []PlayParticipant{
// 			{
// 				Participant: Participant{
// 					ID:        "ID1",
// 					Name:      "name1",
// 					Phone:     "phone1",
// 					Note:      "note1",
// 					CreatedAt: s.mockTime,
// 				},
// 				TotalDonation:      300,
// 				TotalTicketsNumber: 10,
// 				Donations: []Donation{
// 					{
// 						ID:            "dID1",
// 						ParticipantID: "id1",
// 						Amount:        300,
// 						CreatedAt:     time.Time{},
// 					},
// 				},
// 			},
// 		},

// 		PlayParticipants: []PlayParticipant{
// 			{
// 				Participant: Participant{
// 					ID:        "ID2",
// 					Name:      "name2",
// 					Phone:     "phone2",
// 					Note:      "note2",
// 					CreatedAt: s.mockTime,
// 				},
// 				TotalDonation:      200,
// 				TotalTicketsNumber: 5,
// 				Donations: []Donation{
// 					{
// 						ID:            "dID2",
// 						ParticipantID: "ID2",
// 						Amount:        200,
// 						CreatedAt:     s.mockTime,
// 					},
// 				},
// 			},
// 			{
// 				Participant: Participant{
// 					ID:        "ID3",
// 					Name:      "name3",
// 					Phone:     "phone3",
// 					Note:      "note3",
// 					CreatedAt: s.mockTime,
// 				},
// 				TotalDonation:      100,
// 				TotalTicketsNumber: 2,
// 				Donations: []Donation{
// 					{
// 						ID:            "dID3",
// 						ParticipantID: "ID3",
// 						Amount:        1000,
// 						CreatedAt:     s.mockTime,
// 					},
// 				},
// 			},
// 		},
// 	}

// 	mockedPrizeID := "pz1"

// 	s.prizeStorage.EXPECT().Get(mockedPrizeID).Return(&przs[0], nil)

// 	res, err := s.manager.PlayPrizeAgain(s.mockUUID, mockedPrizeID, mockedPreviousResult)
// 	s.Require().NoError(err)
// 	s.Require().NotNil(res)
// 	s.Require().NotEmpty(res.Winners)
// 	s.Require().NotEmpty(res.PlayParticipants)
// }

func setUUIDMock(uuid string) {
	stringUUID = func() string {
		return uuid
	}
}

func setTimeNowMock(t time.Time) {
	timeNow = func() time.Time {
		return t
	}
}

func dummyRaffleRequest() *RaffleRequest {
	return &RaffleRequest{
		Name: "Note name",
		Note: "Note test",
	}
}

func dummyRaffle() *Raffle {
	return &Raffle{
		ID:        "raffle_id",
		Name:      "raffleName",
		Note:      "raffle note",
		CreatedAt: timeNow(),
	}
}
