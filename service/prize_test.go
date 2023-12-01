package service

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type PrizeSuite struct {
	suite.Suite

	ctrl     *gomock.Controller
	storage  *MockPrizeStorage
	manager  *PrizeManager
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
	s.manager = NewPrizeManager(s.storage)
}

func (s *PrizeSuite) TestCreatePrize() {
	prizeRequest := dummyPrizeRequest()

	mockedPrize := Prize{
		ID:          s.mockUUID,
		Name:        prizeRequest.Name,
		TicketCost:  prizeRequest.TicketCost,
		Description: prizeRequest.Description,
		CreatedAt:   s.mockTime,
	}

	s.storage.EXPECT().Create(&mockedPrize).Return(nil)

	resID, err := s.manager.Create(prizeRequest)
	require.NoError(s.T(), err)
	require.Equal(s.T(), mockedPrize.ID, resID)

	s.Run("error", func() {
		request := dummyPrizeRequest()
		expectedPrize := &Prize{
			ID:          s.mockUUID,
			Name:        request.Name,
			TicketCost:  request.TicketCost,
			Description: request.Description,
			CreatedAt:   s.mockTime,
		}

		s.storage.EXPECT().Create(expectedPrize).Return(assert.AnError)

		resID, err := s.manager.Create(request)
		require.ErrorIs(s.T(), err, assert.AnError)
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
	mockedPrize := dummyPrize()
	prizeRequest := dummyPrizeRequest()

	s.storage.EXPECT().Get(mockedPrize.ID).Return(mockedPrize, nil)
	s.storage.EXPECT().Update(mockedPrize).Return(nil)

	err := s.manager.Edit(mockedPrize.ID, prizeRequest)
	require.NoError(s.T(), err)

	s.Run("error", func() {
		s.storage.EXPECT().Get(mockedPrize.ID).Return(nil, assert.AnError)

		err := s.manager.Edit(mockedPrize.ID, prizeRequest)
		require.ErrorIs(s.T(), err, assert.AnError)
	})
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
