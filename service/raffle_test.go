package service

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRaffle(t *testing.T) {
	ctrl := gomock.NewController(t)

	rsMock := NewMockRaffleStorage(ctrl)
	rm := NewRaffleManager(rsMock)

	req := RaffleRequest{
		Name: "raffle_name_1",
		Note: "raffle_note_1",
	}

	mockedID := "raffle_id_1"
	mockedTime := time.Now().UTC()
	mockedErr := assert.AnError

	mockedRaffle := Raffle{
		ID:        mockedID,
		Name:      req.Name,
		Note:      req.Note,
		CreatedAt: mockedTime,
	}

	t.Run("create", func(t *testing.T) {
		t.Run("error", func(t *testing.T) {
			rsMock.EXPECT().Create(gomock.Any()).Return(mockedErr).Times(1)

			res, err := rm.Create(&req)
			require.ErrorIs(t, err, mockedErr)
			require.Empty(t, res)
		})

		t.Run("success", func(t *testing.T) {
			setUUIDMock(mockedID)
			setTimeNowMock(mockedTime)

			rsMock.EXPECT().Create(&mockedRaffle).Return(nil).Times(1)

			resID, err := rm.Create(&req)
			require.NoError(t, err)
			require.Equal(t, mockedID, resID)
		})
	})

	t.Run("get", func(t *testing.T) {
		t.Run("error", func(t *testing.T) {
			rsMock.EXPECT().Get(mockedID).Return(nil, mockedErr).Times(1)

			res, err := rm.Get(mockedID)
			require.ErrorIs(t, err, mockedErr)
			require.Nil(t, res)
		})

		t.Run("success", func(t *testing.T) {
			rsMock.EXPECT().Get(mockedID).Return(&mockedRaffle, nil).Times(1)

			id, err := rm.Get(mockedID)
			require.NoError(t, err)
			require.Equal(t, &mockedRaffle, id)
		})
	})

	t.Run("list", func(t *testing.T) {
		t.Run("error", func(t *testing.T) {
			rsMock.EXPECT().GetAll().Return(nil, mockedErr).Times(1)

			res, err := rm.List()
			require.ErrorIs(t, err, mockedErr)
			require.Nil(t, res)
		})

		t.Run("success", func(t *testing.T) {
			raffles := []Raffle{mockedRaffle, mockedRaffle}

			rsMock.EXPECT().GetAll().Return(raffles, nil).Times(1)

			res, err := rm.List()
			require.NoError(t, err)
			require.Equal(t, raffles, res)
		})
	})

	t.Run("Export non-empty collection s", func(t *testing.T) {
		raf := &Raffle{ID: mockedID, Name: "Raffle Test"}
		prts := []Participant{
			{ID: "p1", Name: "Participant 1"},
			{ID: "p2", Name: "Participant 2"},
		}
		przs := []Prize{
			{ID: "pr1", Name: "Prize 1"},
			{ID: "pr2", Name: "Prize 2"},
		}

		rsMock.EXPECT().Get(mockedID).Return(raf, nil).Times(1)

		psMock := NewMockParticipantStorage(ctrl)
		rsMock.EXPECT().ParticipantStorage(mockedID).Return(psMock).Times(1)
		psMock.EXPECT().GetAll().Return(prts, nil).Times(1)

		pzMock := NewMockPrizeStorage(ctrl)
		rsMock.EXPECT().PrizeStorage(mockedID).Return(pzMock).Times(1)
		pzMock.EXPECT().GetAll().Return(przs, nil).Times(1)

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
