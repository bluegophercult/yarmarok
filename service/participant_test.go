package service_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaznasho/yarmarok/mocks"
	"github.com/kaznasho/yarmarok/service"
	"github.com/stretchr/testify/assert"
)

func TestParticipantManagerAdd(t *testing.T) {
	ctrl := gomock.NewController(t)

	storageMock := mocks.NewMockParticipantStorage(ctrl)
	manager := service.NewParticipantManager(storageMock)

	t.Run("add participant", func(t *testing.T) {
		storageMock.EXPECT().Create(gomock.Any()).Return(nil)

		_, err := manager.Create(&service.ParticipantAddRequest{
			Name:  "John Doe",
			Phone: "1234567890",
			Note:  "Test participant",
		})

		assert.NoError(t, err)
	})

	t.Run("add_already_exists", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		storageMock.EXPECT().Create(gomock.Any()).Return(service.ErrParticipantAlreadyExists)

		participantManager := service.NewParticipantManager(storageMock)

		_, err := participantManager.Create(&service.ParticipantAddRequest{
			Name:  "John Doe",
			Phone: "1234567890",
			Note:  "Test participant",
		})

		assert.ErrorIs(t, err, service.ErrParticipantAlreadyExists)
	})
}

func TestParticipantManagerEdit(t *testing.T) {
	ctrl := gomock.NewController(t)

	storageMock := mocks.NewMockParticipantStorage(ctrl)
	manager := service.NewParticipantManager(storageMock)

	t.Run("edit participant", func(t *testing.T) {
		storageMock.EXPECT().Get(gomock.Any()).Return(&service.Participant{}, nil)
		storageMock.EXPECT().Update(gomock.Any()).Return(nil)

		_, err := manager.Edit(&service.ParticipantEditRequest{ID: "test-id"})

		assert.NoError(t, err)
	})

	t.Run("participant not found", func(t *testing.T) {
		storageMock.EXPECT().Get(gomock.Any()).Return(nil, service.ErrParticipantNotFound)

		_, err := manager.Edit(&service.ParticipantEditRequest{
			ID: "test-id",
		})

		assert.ErrorIs(t, err, service.ErrParticipantNotFound)
	})
}

func TestParticipantManagerList(t *testing.T) {
	ctrl := gomock.NewController(t)

	storageMock := mocks.NewMockParticipantStorage(ctrl)
	manager := service.NewParticipantManager(storageMock)

	t.Run("list participants", func(t *testing.T) {
		storageMock.EXPECT().GetAll().Return([]service.Participant{}, nil)

		_, err := manager.List()

		assert.NoError(t, err)
	})
}
