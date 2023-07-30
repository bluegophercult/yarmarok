package service

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

//go:generate mockgen -destination=mock_participant_storage_test.go -package=service github.com/kaznasho/yarmarok/service ParticipantStorage

func TestParticipantManagerAdd(t *testing.T) {
	ctrl := gomock.NewController(t)

	storageMock := NewMockParticipantStorage(ctrl)
	manager := NewParticipantManager(storageMock)

	t.Run("add participant", func(t *testing.T) {
		storageMock.EXPECT().Create(gomock.Any()).Return(nil)

		_, err := manager.Create(&ParticipantRequest{
			Name:  "John Doe",
			Phone: "1234567890",
			Note:  "Test participant",
		})

		assert.NoError(t, err)
	})

	t.Run("add_already_exists", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		storageMock.EXPECT().Create(gomock.Any()).Return(ErrParticipantAlreadyExists)

		participantManager := NewParticipantManager(storageMock)

		_, err := participantManager.Create(&ParticipantRequest{
			Name:  "John Doe",
			Phone: "1234567890",
			Note:  "Test participant",
		})

		assert.ErrorIs(t, err, ErrParticipantAlreadyExists)
	})
}

func TestParticipantManagerEdit(t *testing.T) {
	ctrl := gomock.NewController(t)

	storageMock := NewMockParticipantStorage(ctrl)
	manager := NewParticipantManager(storageMock)

	t.Run("edit participant", func(t *testing.T) {
		storageMock.EXPECT().Get(gomock.Any()).Return(&Participant{}, nil)
		storageMock.EXPECT().Update(gomock.Any()).Return(nil)

		_, err := manager.Edit(&ParticipantEditRequest{ID: "test-id"})

		assert.NoError(t, err)
	})

	t.Run("participant not found", func(t *testing.T) {
		storageMock.EXPECT().Get(gomock.Any()).Return(nil, ErrParticipantNotFound)

		_, err := manager.Edit(&ParticipantEditRequest{
			ID: "test-id",
		})

		assert.ErrorIs(t, err, ErrParticipantNotFound)
	})
}

func TestParticipantManagerList(t *testing.T) {
	ctrl := gomock.NewController(t)

	storageMock := NewMockParticipantStorage(ctrl)
	manager := NewParticipantManager(storageMock)

	t.Run("list participants", func(t *testing.T) {
		storageMock.EXPECT().GetAll().Return([]Participant{}, nil)

		_, err := manager.List()

		assert.NoError(t, err)
	})
}
