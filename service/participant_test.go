package service

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:generate mockgen -destination=mock_participant_storage_test.go -package=service github.com/kaznasho/yarmarok/service ParticipantStorage

func TestParticipantManagerAdd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParticipantStorage := NewMockParticipantStorage(ctrl)

	participantManager := NewParticipantManager(mockParticipantStorage)

	t.Run("add participant", func(t *testing.T) {
		mockParticipantStorage.EXPECT().Create(gomock.Any()).Return(nil)

		_, err := participantManager.Add(&ParticipantInitRequest{
			Name:  "John Doe",
			Phone: "1234567890",
			Note:  "Test participant",
		})

		assert.NoError(t, err)
	})

	t.Run("participant already exists", func(t *testing.T) {
		mockParticipantStorage.EXPECT().Create(gomock.Any()).Return(ErrParticipantAlreadyExists)

		_, err := participantManager.Add(&ParticipantInitRequest{
			Name:  "John Doe",
			Phone: "1234567890",
			Note:  "Test participant",
		})

		assert.ErrorIs(t, err, ErrParticipantAlreadyExists)
	})
}

func TestParticipantManagerEdit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParticipantStorage := NewMockParticipantStorage(ctrl)

	participantManager := NewParticipantManager(mockParticipantStorage)

	t.Run("edit participant", func(t *testing.T) {
		mockParticipantStorage.EXPECT().Get(gomock.Any()).Return(&Participant{}, nil)
		mockParticipantStorage.EXPECT().Update(gomock.Any()).Return(nil)

		_, err := participantManager.Edit(&ParticipantEditRequest{
			ID: "test-id",
		})

		assert.NoError(t, err)
	})

	t.Run("participant not found", func(t *testing.T) {
		mockParticipantStorage.EXPECT().Get(gomock.Any()).Return(nil, ErrParticipantNotFound)

		_, err := participantManager.Edit(&ParticipantEditRequest{
			ID: "test-id",
		})

		assert.ErrorIs(t, err, ErrParticipantNotFound)
	})
}

func TestParticipantManagerList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParticipantStorage := NewMockParticipantStorage(ctrl)

	participantManager := NewParticipantManager(mockParticipantStorage)

	t.Run("list participants", func(t *testing.T) {
		mockParticipantStorage.EXPECT().GetAll().Return([]Participant{}, nil)

		_, err := participantManager.List()

		assert.NoError(t, err)
	})
}
