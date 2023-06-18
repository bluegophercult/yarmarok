package service

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

//go:generate mockgen -destination=mock_organizer_storage_test.go -package=service github.com/kaznasho/yarmarok/service OrganizerStorage
func TestInitOrganizer(t *testing.T) {
	ctrl := gomock.NewController(t)

	storageMock := NewMockOrganizerStorage(ctrl)

	t.Run("init organizer", func(t *testing.T) {
		t.Run("exists", func(t *testing.T) {
			organizerID := "123"
			storageMock.EXPECT().Exists(organizerID).Return(true, nil)
			manager := NewOrganizerManager(storageMock)

			err := manager.InitOrganizerIfNotExists(organizerID)
			assert.NoError(t, err)
		})

		t.Run("not exists", func(t *testing.T) {
			organizerID := "123"
			storageMock.EXPECT().Exists(organizerID).Return(false, nil)
			storageMock.EXPECT().Create(Organizer{ID: organizerID}).Return(nil)
			manager := NewOrganizerManager(storageMock)

			err := manager.InitOrganizerIfNotExists(organizerID)
			assert.NoError(t, err)
		})

		t.Run("error", func(t *testing.T) {
			organizerID := "123"
			storageMock.EXPECT().Exists(organizerID).Return(false, assert.AnError)
			manager := NewOrganizerManager(storageMock)

			err := manager.InitOrganizerIfNotExists(organizerID)
			assert.Error(t, err)
		})
	})
}

var _ OrganizerService = &OrganizerManager{}
