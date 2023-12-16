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
