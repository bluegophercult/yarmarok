package service

import (
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

//go:generate mockgen -destination=mock_yarmarok_storage_test.go -package=service github.com/kaznasho/yarmarok/service YarmarokStorage

func TestYarmarok(t *testing.T) {
	ctrl := gomock.NewController(t)

	storageMock := NewMockYarmarokStorage(ctrl)

	manager := NewYarmarokManager(storageMock)

	t.Run("init", func(t *testing.T) {
		t.Run("error", func(t *testing.T) {
			req := YarmarokInitRequest{
				Name: "yarmarok_name_1",
				Note: "yarmarok_note_1",
			}

			mockedErr := assert.AnError

			storageMock.EXPECT().Create(gomock.Any()).Return(mockedErr).Times(1)

			res, err := manager.Init(&req)
			assert.Error(t, err)
			assert.Equal(t, mockedErr, err)
			assert.Nil(t, res)
		})

		t.Run("success", func(t *testing.T) {
			req := YarmarokInitRequest{
				Name: "yarmarok_name_1",
				Note: "yarmarok_note_1",
			}

			mockedID := "yarmarok_id_1"
			mockedTime := time.Now().UTC()

			setUUIDMock(mockedID)
			setTimeNowMock(mockedTime)

			mockedYarmarok := &Yarmarok{
				ID:        mockedID,
				Name:      req.Name,
				Note:      req.Note,
				CreatedAt: mockedTime,
			}

			storageMock.EXPECT().Create(mockedYarmarok).Return(nil).Times(1)

			res, err := manager.Init(&req)
			assert.NoError(t, err)
			assert.Equal(t, mockedID, res.ID)
		})
	})

	t.Run("get", func(t *testing.T) {
		t.Run("error", func(t *testing.T) {
			id := "yarmarok_id_1"

			mockedErr := assert.AnError

			storageMock.EXPECT().Get(id).Return(nil, mockedErr).Times(1)

			res, err := manager.Get(id)
			assert.Error(t, err)
			assert.Equal(t, mockedErr, err)
			assert.Nil(t, res)
		})

		t.Run("success", func(t *testing.T) {
			id := "yarmarok_id_1"

			mockedYarmarok := &Yarmarok{
				ID:          id,
				Name:        "yarmarok_name_1",
				Note:        "yarmarok_note_1",
				CreatedAt:   timeNow().UTC(),
				OrganizerID: "organizer_id_1",
			}

			storageMock.EXPECT().Get(id).Return(mockedYarmarok, nil).Times(1)

			res, err := manager.Get(id)
			assert.NoError(t, err)
			assert.Equal(t, mockedYarmarok, res)
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
			mockedYarmaroks := []Yarmarok{
				{
					ID:          "yarmarok_id_1",
					Name:        "yarmarok_name_1",
					Note:        "yarmarok_note_1",
					CreatedAt:   timeNow().UTC(),
					OrganizerID: "organizer_id_1",
				},
				{
					ID:          "yarmarok_id_2",
					Name:        "yarmarok_name_2",
					Note:        "yarmarok_note_2",
					CreatedAt:   timeNow().UTC(),
					OrganizerID: "organizer_id_1",
				},
			}

			expected := &YarmarokListResponse{
				Yarmaroks: mockedYarmaroks,
			}

			storageMock.EXPECT().GetAll().Return(mockedYarmaroks, nil).Times(1)

			res, err := manager.List()
			assert.NoError(t, err)
			assert.Equal(t, expected, res)
		})
	})
}

var _ YarmarokService = &YarmarokManager{}

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
