package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestInitOrganizer(t *testing.T) {
	ctrl := gomock.NewController(t)

	osMock := NewMockOrganizerStorage(ctrl)

	t.Run("init organizer", func(t *testing.T) {
		t.Run("exists", func(t *testing.T) {
			organizerID := "123"
			osMock.EXPECT().Exists(organizerID).Return(true, nil)
			om := NewOrganizerManager(osMock)

			err := om.CreateOrganizerIfNotExists(organizerID)
			assert.NoError(t, err)
		})

		t.Run("not exists", func(t *testing.T) {
			organizerID := "123"
			osMock.EXPECT().Exists(organizerID).Return(false, nil)
			osMock.EXPECT().Create(&Organizer{ID: organizerID}).Return(nil)
			om := NewOrganizerManager(osMock)

			err := om.CreateOrganizerIfNotExists(organizerID)
			assert.NoError(t, err)
		})

		t.Run("error", func(t *testing.T) {
			organizerID := "123"
			osMock.EXPECT().Exists(organizerID).Return(false, assert.AnError)
			om := NewOrganizerManager(osMock)

			err := om.CreateOrganizerIfNotExists(organizerID)
			assert.Error(t, err)
		})
	})
}
