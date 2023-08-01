package service_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaznasho/yarmarok/mocks"
	"github.com/kaznasho/yarmarok/service"
	"github.com/stretchr/testify/assert"
)

func TestPrizeManager(t *testing.T) {
	ctrl := gomock.NewController(t)

	storageMock := mocks.NewMockPrizeStorage(ctrl)
	manager := service.NewPrizeManager(storageMock)

	testID := "prize_id_1"

	t.Run("add", func(t *testing.T) {
		t.Run("prize_success", func(t *testing.T) {
			storageMock.EXPECT().Create(gomock.Any()).Return(nil)

			res, err := manager.Create(&service.PrizeCreateRequest{
				Name:        "prize_name_1",
				TicketCost:  1234,
				Description: "prize_description_1",
			})

			assert.NoError(t, err)
			assert.NotNil(t, res)
		})

		t.Run("prize_already_exists", func(t *testing.T) {
			storageMock.EXPECT().Create(gomock.Any()).Return(service.ErrAlreadyExists)

			prizeManager := service.NewPrizeManager(storageMock)

			res, err := prizeManager.Create(&service.PrizeCreateRequest{
				Name:        "prize_name_1",
				TicketCost:  1234,
				Description: "prize_description_1",
			})

			assert.Error(t, err)
			assert.ErrorIs(t, err, service.ErrAlreadyExists)
			assert.Nil(t, res)
		})
	})

	t.Run("edit", func(t *testing.T) {
		t.Run("prize_success", func(t *testing.T) {
			storageMock.EXPECT().Get(gomock.Any()).Return(&service.Prize{}, nil)
			storageMock.EXPECT().Update(gomock.Any()).Return(nil)

			res, err := manager.Edit(&service.PrizeEditRequest{ID: testID})

			assert.NoError(t, err)
			assert.NotNil(t, res)
		})

		t.Run("prize_not_found", func(t *testing.T) {
			storageMock.EXPECT().Get(gomock.Any()).Return(nil, service.ErrAlreadyExists)

			res, err := manager.Edit(&service.PrizeEditRequest{ID: testID})

			assert.Error(t, err)
			assert.ErrorIs(t, err, service.ErrAlreadyExists)
			assert.Equal(t, &service.Result{service.StatusError}, res)
		})

		t.Run("update_prize_error", func(t *testing.T) {
			mockedErr := assert.AnError

			storageMock.EXPECT().Get(gomock.Any()).Return(&service.Prize{}, nil)
			storageMock.EXPECT().Update(gomock.Any()).Return(mockedErr)

			res, err := manager.Edit(&service.PrizeEditRequest{ID: testID})

			assert.Error(t, err)
			assert.ErrorIs(t, err, mockedErr)
			assert.Equal(t, &service.Result{service.StatusError}, res)
		})
	})

	t.Run("list", func(t *testing.T) {

		t.Run("prize_success", func(t *testing.T) {
			storageMock.EXPECT().GetAll().Return([]service.Prize{}, nil)

			res, err := manager.List()

			assert.NoError(t, err)
			assert.NotNil(t, res)
		})

		t.Run("prize_error", func(t *testing.T) {
			mockedErr := assert.AnError

			storageMock.EXPECT().GetAll().Return(nil, mockedErr)

			res, err := manager.List()

			assert.Error(t, err)
			assert.ErrorIs(t, err, mockedErr)
			assert.Nil(t, res)
		})
	})
}
