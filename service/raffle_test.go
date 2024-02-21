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
	s.storage.EXPECT().ParticipantStorage(s.mockUUID).Return(psMock).Times(2)
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
