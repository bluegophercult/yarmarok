package service_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaznasho/yarmarok/mocks"
	"github.com/kaznasho/yarmarok/service"
	"github.com/stretchr/testify/assert"
)

func TestInitOrganizer(t *testing.T) {
	ctrl := gomock.NewController(t)

	osMock := mocks.NewMockOrganizerStorage(ctrl)

	t.Run("init organizer", func(t *testing.T) {
		t.Run("exists", func(t *testing.T) {
			organizerID := "123"
			osMock.EXPECT().Exists(organizerID).Return(true, nil)
			om := service.NewOrganizerManager(osMock)

			err := om.CreateOrganizerIfNotExists(organizerID)
			assert.NoError(t, err)
		})

		t.Run("not exists", func(t *testing.T) {
			organizerID := "123"
			osMock.EXPECT().Exists(organizerID).Return(false, nil)
			osMock.EXPECT().Create(&service.Organizer{ID: organizerID}).Return(nil)
			om := service.NewOrganizerManager(osMock)

			err := om.CreateOrganizerIfNotExists(organizerID)
			assert.NoError(t, err)
		})

		t.Run("error", func(t *testing.T) {
			organizerID := "123"
			osMock.EXPECT().Exists(organizerID).Return(false, assert.AnError)
			om := service.NewOrganizerManager(osMock)

			err := om.CreateOrganizerIfNotExists(organizerID)
			assert.Error(t, err)
		})
	})
}
