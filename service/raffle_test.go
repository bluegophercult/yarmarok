package service

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RaffleSuite struct {
	suite.Suite

	ctrl    *gomock.Controller
	storage *MockRaffleStorage
	manager *RaffleManager
}

func TestRaffle(t *testing.T) {
	suite.Run(t, &RaffleSuite{})
}

func (s *RaffleSuite) SetupTest() {
	setTimeNowMock(time.Now().UTC())
	setUUIDMock("raffle_id_1")

	s.ctrl = gomock.NewController(s.T())
	s.storage = NewMockRaffleStorage(s.ctrl)
	s.manager = NewRaffleManager(s.storage)
}

func (s *RaffleSuite) TestCreateRaffle() {
	raffleRequest := &RaffleRequest{
		Name: "Example",
		Note: "eeee",
	}

	mockedRaffle := Raffle{
		ID:        "raffle_id_1",
		Name:      raffleRequest.Name,
		Note:      raffleRequest.Note,
		CreatedAt: timeNow(),
	}

	s.storage.EXPECT().Create(&mockedRaffle).Return(nil)

	resID, err := s.manager.Create(raffleRequest)
	require.NoError(s.T(), err)
	require.Equal(s.T(), mockedRaffle.ID, resID)
}

func TestRafflex(t *testing.T) {
	ctrl := gomock.NewController(t)

	storageMock := NewMockRaffleStorage(ctrl)
	rm := NewRaffleManager(storageMock)

	raffleRequest := RaffleRequest{
		Name: "raffle_name_1",
		Note: "raffle_note_1",
	}

	mockedID := "raffle_id_1"
	mockedTime := time.Now().UTC()
	mockedErr := assert.AnError

	mockedRaffle := Raffle{
		ID:        mockedID,
		Name:      raffleRequest.Name,
		Note:      raffleRequest.Note,
		CreatedAt: mockedTime,
	}

	setTimeNowMock(mockedTime)
	setUUIDMock(mockedID)

	t.Run("create", func(t *testing.T) {
		t.Run("error", func(t *testing.T) {
			request := dummyRaffleRequest()
			expectedRaffle := &Raffle{
				ID:        mockedID,
				Name:      request.Name,
				Note:      request.Note,
				CreatedAt: mockedTime,
			}

			storageMock.EXPECT().Create(expectedRaffle).Return(mockedErr)

			response, err := rm.Create(request)
			assert.ErrorIs(t, err, mockedErr)
			assert.Equal(t, "", response)
		})

		t.Run("success", func(t *testing.T) {
			request := dummyRaffleRequest()
			expectedRaffle := &Raffle{
				ID:        mockedID,
				Name:      request.Name,
				Note:      request.Note,
				CreatedAt: mockedTime,
			}

			storageMock.EXPECT().Create(expectedRaffle).Return(nil)

			resID, err := rm.Create(request)
			assert.NoError(t, err)
			assert.Equal(t, expectedRaffle.ID, resID)
		})

		t.Run("invalid_name", func(t *testing.T) {
			request := dummyRaffleRequest()
			request.Name = ""

			response, err := rm.Create(request)
			assert.ErrorIs(t, err, ErrInvalidRequest)
			assert.Equal(t, "", response)
		})

		t.Run("invalid_note", func(t *testing.T) {
			request := dummyRaffleRequest()
			request.Note = "<>////"

			response, err := rm.Create(request)
			assert.ErrorIs(t, err, ErrInvalidRequest)
			assert.Equal(t, "", response)
		})
	})

	t.Run("get", func(t *testing.T) {
		t.Run("error", func(t *testing.T) {
			storageMock.EXPECT().Get(mockedID).Return(nil, mockedErr)

			res, err := rm.Get(mockedID)
			require.ErrorIs(t, err, mockedErr)
			require.Nil(t, res)
		})

		t.Run("success", func(t *testing.T) {
			mockedRaffle := dummyRaffle()
			storageMock.EXPECT().Get(mockedID).Return(mockedRaffle, nil)

			raf, err := rm.Get(mockedID)
			require.NoError(t, err)
			require.Equal(t, mockedRaffle, raf)
		})
	})

	t.Run("edit", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			raffleRequest := dummyRaffleRequest()
			mockedRaffle := Raffle{
				ID:        "raffle_id_1",
				Name:      raffleRequest.Name,
				Note:      raffleRequest.Note,
				CreatedAt: mockedTime,
			}

			storageMock.EXPECT().Get(mockedID).Return(&mockedRaffle, nil)
			storageMock.EXPECT().Update(&mockedRaffle).Return(nil)
			err := rm.Edit(mockedID, raffleRequest)
			require.NoError(t, err)
		})

		t.Run("not_found", func(t *testing.T) {
			storageMock.EXPECT().Get(mockedID).Return(nil, ErrNotFound)
			err := rm.Edit(mockedID, &raffleRequest)
			require.ErrorIs(t, err, ErrNotFound)
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			storageMock.EXPECT().Delete(mockedID).Return(nil)
			err := rm.Delete(mockedID)
			require.NoError(t, err)
		})

		t.Run("not_found", func(t *testing.T) {
			storageMock.EXPECT().Delete(mockedID).Return(ErrNotFound)
			err := rm.Delete(mockedID)
			require.ErrorIs(t, err, ErrNotFound)
		})
	})

	t.Run("list", func(t *testing.T) {
		t.Run("error", func(t *testing.T) {
			storageMock.EXPECT().GetAll().Return(nil, mockedErr)

			res, err := rm.List()
			require.ErrorIs(t, err, mockedErr)
			require.Nil(t, res)
		})

		t.Run("success", func(t *testing.T) {
			raffles := []Raffle{mockedRaffle, mockedRaffle, mockedRaffle}

			storageMock.EXPECT().GetAll().Return(raffles, nil)

			res, err := rm.List()
			require.NoError(t, err)
			require.Equal(t, raffles, res)
		})
	})

	t.Run("Export non-empty collections", func(t *testing.T) {
		raf := &Raffle{ID: mockedID, Name: "Raffle Test"}
		prts := []Participant{
			{ID: "p1", Name: "Participant 1"},
			{ID: "p2", Name: "Participant 2"},
		}
		przs := []Prize{
			{ID: "pr1", Name: "Prize 1"},
			{ID: "pr2", Name: "Prize 2"},
		}

		storageMock.EXPECT().Get(mockedID).Return(raf, nil)

		psMock := NewMockParticipantStorage(ctrl)
		storageMock.EXPECT().ParticipantStorage(mockedID).Return(psMock)
		psMock.EXPECT().GetAll().Return(prts, nil)

		pzMock := NewMockPrizeStorage(ctrl)
		storageMock.EXPECT().PrizeStorage(mockedID).Return(pzMock)
		pzMock.EXPECT().GetAll().Return(przs, nil)

		res, err := rm.Export(mockedID)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, "yarmarok_"+mockedID+".xlsx", res.FileName)
		require.NotEmpty(t, res.Content)
	})
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
