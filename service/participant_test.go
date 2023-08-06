package service

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

//go:generate mockgen -destination=mock_participant_storage_test.go -package=service github.com/kaznasho/yarmarok/service ParticipantStorage

func TestParticipantManagerAdd(t *testing.T) {
	ctrl := gomock.NewController(t)

	storageMock := NewMockParticipantStorage(ctrl)
	manager := NewParticipantManager(storageMock)

	prt := &ParticipantRequest{
		Name:  "John Doe",
		Phone: "1234567890",
		Note:  "Test participant",
	}

	t.Run("add participant", func(t *testing.T) {
		storageMock.EXPECT().Create(gomock.Any()).Return(nil)
		_, err := manager.Create(prt)
		require.NoError(t, err)
	})

	t.Run("add_already_exists", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		storageMock.EXPECT().Create(gomock.Any()).Return(ErrParticipantAlreadyExists)

		participantManager := NewParticipantManager(storageMock)
		_, err := participantManager.Create(prt)
		require.ErrorIs(t, err, ErrParticipantAlreadyExists)
	})
}

func TestParticipantManagerEdit(t *testing.T) {
	ctrl := gomock.NewController(t)

	storageMock := NewMockParticipantStorage(ctrl)
	manager := NewParticipantManager(storageMock)

	id := "participant_id_1"
	prt := &ParticipantRequest{Name: "John Doe", Phone: "1234567890", Note: "Test participant"}
	t.Run("edit participant", func(t *testing.T) {
		storageMock.EXPECT().Get(gomock.Any()).Return(&Participant{}, nil)
		storageMock.EXPECT().Update(gomock.Any()).Return(nil)

		err := manager.Edit(id, prt)
		require.NoError(t, err)
	})

	t.Run("participant not found", func(t *testing.T) {
		storageMock.EXPECT().Get(gomock.Any()).Return(nil, ErrParticipantNotFound)

		err := manager.Edit(id, prt)
		require.ErrorIs(t, err, ErrParticipantNotFound)
	})
}

func TestParticipantManagerList(t *testing.T) {
	ctrl := gomock.NewController(t)

	storageMock := NewMockParticipantStorage(ctrl)
	manager := NewParticipantManager(storageMock)

	t.Run("list participants", func(t *testing.T) {
		storageMock.EXPECT().GetAll().Return([]Participant{}, nil)

		_, err := manager.List()

		require.NoError(t, err)
	})
}
