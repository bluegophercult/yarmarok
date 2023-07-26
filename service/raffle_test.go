package service

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:generate mockgen -destination=mock_raffle_storage_test.go -package=service github.com/kaznasho/yarmarok/service RaffleStorage

func TestRaffle(t *testing.T) {
	ctrl := gomock.NewController(t)

	storageMock := NewMockRaffleStorage(ctrl)

	manager := NewRaffleManager(storageMock)

	t.Run("init", func(t *testing.T) {
		t.Run("error", func(t *testing.T) {
			req := RaffleInitRequest{
				Name: "raffle_name_1",
				Note: "raffle_note_1",
			}

			mockedErr := assert.AnError

			storageMock.EXPECT().Create(gomock.Any()).Return(mockedErr).Times(1)

			res, err := manager.Create(&req)
			assert.Error(t, err)
			assert.Equal(t, mockedErr, err)
			assert.Nil(t, res)
		})

		t.Run("success", func(t *testing.T) {
			req := RaffleInitRequest{
				Name: "raffle_name_1",
				Note: "raffle_note_1",
			}

			mockedID := "raffle_id_1"
			mockedTime := time.Now().UTC()

			setUUIDMock(mockedID)
			setTimeNowMock(mockedTime)

			mockedRaffle := &Raffle{
				ID:        mockedID,
				Name:      req.Name,
				Note:      req.Note,
				CreatedAt: mockedTime,
			}

			storageMock.EXPECT().Create(mockedRaffle).Return(nil).Times(1)

			res, err := manager.Create(&req)
			assert.NoError(t, err)
			assert.Equal(t, mockedID, res.ID)
		})
	})

	t.Run("get", func(t *testing.T) {
		t.Run("error", func(t *testing.T) {
			id := "raffle_id_1"

			mockedErr := assert.AnError

			storageMock.EXPECT().Get(id).Return(nil, mockedErr).Times(1)

			res, err := manager.Get(id)
			assert.Error(t, err)
			assert.Equal(t, mockedErr, err)
			assert.Nil(t, res)
		})

		t.Run("success", func(t *testing.T) {
			id := "raffle_id_1"

			mockedRaffle := &Raffle{
				ID:          id,
				Name:        "raffle_name_1",
				Note:        "raffle_note_1",
				CreatedAt:   timeNow().UTC(),
				OrganizerID: "organizer_id_1",
			}

			storageMock.EXPECT().Get(id).Return(mockedRaffle, nil).Times(1)

			res, err := manager.Get(id)
			assert.NoError(t, err)
			assert.Equal(t, mockedRaffle, res)
		})
	})

	t.Run("list", func(t *testing.T) {
		t.Run("error", func(t *testing.T) {
			mockedErr := assert.AnError

			storageMock.EXPECT().GetAll().Return(nil, mockedErr).Times(1)

			res, err := manager.List()
			assert.Error(t, err)
			assert.ErrorIs(t, err, mockedErr)
			assert.Nil(t, res)
		})

		t.Run("success", func(t *testing.T) {
			mockedRaffles := []Raffle{
				{
					ID:          "raffle_id_1",
					Name:        "raffle_name_1",
					Note:        "raffle_note_1",
					CreatedAt:   timeNow().UTC(),
					OrganizerID: "organizer_id_1",
				},
				{
					ID:          "raffle_id_2",
					Name:        "raffle_name_2",
					Note:        "raffle_note_2",
					CreatedAt:   timeNow().UTC(),
					OrganizerID: "organizer_id_1",
				},
			}

			expected := &RaffleListResponse{
				Raffles: mockedRaffles,
			}

			storageMock.EXPECT().GetAll().Return(mockedRaffles, nil).Times(1)

			res, err := manager.List()
			assert.NoError(t, err)
			assert.Equal(t, expected, res)
		})
	})

	t.Run("Export non-empty collection s", func(t *testing.T) {
		id := "raffle_id_1"
		raf := &Raffle{ID: id, Name: "Raffle Test"} // Add more fields as needed
		participants := []Participant{
			{ID: "p1", Name: "Participant 1"},
			{ID: "p2", Name: "Participant 2"},
		}
		prizes := []Prize{
			{ID: "pr1", Name: "Prize 1"},
			{ID: "pr2", Name: "Prize 2"},
		}

		storageMock.EXPECT().Get(id).Return(raf, nil).Times(1)

		prtStorage := NewMockParticipantStorage(ctrl)
		storageMock.EXPECT().ParticipantStorage(id).Return(prtStorage).Times(1)
		prtStorage.EXPECT().GetAll().Return(participants, nil).Times(1)

		przStorage := NewMockPrizeStorage(ctrl)
		storageMock.EXPECT().PrizeStorage(id).Return(przStorage).Times(1)
		przStorage.EXPECT().GetAll().Return(prizes, nil).Times(1)

		resp, err := manager.Export(id)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "yarmarok_"+id+".xlsx", resp.FileName)
		require.NotEmpty(t, resp.Content)
	})
}

var _ RaffleService = &RaffleManager{}

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
