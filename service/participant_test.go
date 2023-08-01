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

		_, err := manager.Create(&service.ParticipantRequest{
			Name:  "John Doe",
			Phone: "1234567890",
			Note:  "Test participant",
		})

		assert.NoError(t, err)
	})

	t.Run("add_already_exists", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		storageMock.EXPECT().Create(gomock.Any()).Return(service.ErrAlreadyExists)

		participantManager := service.NewParticipantManager(storageMock)

		_, err := participantManager.Create(&service.ParticipantRequest{
			Name:  "John Doe",
			Phone: "1234567890",
			Note:  "Test participant",
		})

		assert.ErrorIs(t, err, service.ErrAlreadyExists)
	})
}

func TestParticipantManagerEdit(t *testing.T) {
	ctrl := gomock.NewController(t)

	storageMock := mocks.NewMockParticipantStorage(ctrl)
	manager := service.NewParticipantManager(storageMock)

	id := "participant_id"
	p := &service.ParticipantRequest{
		Name:  "John Doe",
		Phone: "1234567890",
		Note:  "Test participant",
	}

	t.Run("edit participant", func(t *testing.T) {
		storageMock.EXPECT().Get(gomock.Any()).Return(&service.Participant{}, nil)
		storageMock.EXPECT().Update(gomock.Any()).Return(nil)

		err := manager.Edit(id, p)

		assert.NoError(t, err)
	})

	t.Run("participant not found", func(t *testing.T) {
		storageMock.EXPECT().Get(gomock.Any()).Return(nil, service.ErrNotFound)

		err := manager.Edit(id, p)

		assert.ErrorIs(t, err, service.ErrNotFound)
	})
}
func TestParticipantManagerDelete(t *testing.T) {
	ctrl := gomock.NewController(t)

	storageMock := mocks.NewMockParticipantStorage(ctrl)
	manager := service.NewParticipantManager(storageMock)

	id := "participant_id"

	t.Run("delete participant", func(t *testing.T) {
		storageMock.EXPECT().Delete(gomock.Any()).Return(nil)

		err := manager.Delete(id)
		assert.NoError(t, err)
	})

	t.Run("participant not found", func(t *testing.T) {
		storageMock.EXPECT().Delete(gomock.Any()).Return(service.ErrNotFound)

		err := manager.Delete(id)
		assert.ErrorIs(t, err, service.ErrNotFound)
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
