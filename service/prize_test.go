package service

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

//go:generate mockgen -destination=mock_prize_storage_test.go -package=service github.com/kaznasho/yarmarok/service PrizeStorage

func TestPrizeManager(t *testing.T) {
	ctrl := gomock.NewController(t)

	storageMock := NewMockPrizeStorage(ctrl)
	manager := NewPrizeManager(storageMock)

	testID := "prize_id_1"

	t.Run("add", func(t *testing.T) {
		t.Run("prize_success", func(t *testing.T) {
			storageMock.EXPECT().Create(gomock.Any()).Return(nil)

			res, err := manager.Add(&PrizeAddRequest{
				Name:        "prize_name_1",
				TicketCost:  1234,
				Description: "prize_description_1",
			})

			assert.NoError(t, err)
			assert.NotNil(t, res)
		})

		t.Run("prize_already_exists", func(t *testing.T) {
			storageMock.EXPECT().Create(gomock.Any()).Return(ErrPrizeAlreadyExists)

			prizeManager := NewPrizeManager(storageMock)

			res, err := prizeManager.Add(&PrizeAddRequest{
				Name:        "prize_name_1",
				TicketCost:  1234,
				Description: "prize_description_1",
			})

			assert.Error(t, err)
			assert.ErrorIs(t, err, ErrPrizeAlreadyExists)
			assert.Nil(t, res)
		})
	})

	t.Run("edit", func(t *testing.T) {
		t.Run("prize_success", func(t *testing.T) {
			storageMock.EXPECT().Get(gomock.Any()).Return(&Prize{}, nil)
			storageMock.EXPECT().Update(gomock.Any()).Return(nil)

			res, err := manager.Edit(&PrizeEditRequest{ID: testID})

			assert.NoError(t, err)
			assert.NotNil(t, res)
		})

		t.Run("prize_not_found", func(t *testing.T) {
			storageMock.EXPECT().Get(gomock.Any()).Return(nil, ErrPrizeAlreadyExists)

			res, err := manager.Edit(&PrizeEditRequest{ID: testID})

			assert.Error(t, err)
			assert.ErrorIs(t, err, ErrPrizeAlreadyExists)
			assert.Equal(t, &Result{StatusError}, res)
		})

		t.Run("update_prize_error", func(t *testing.T) {
			mockedErr := assert.AnError

			storageMock.EXPECT().Get(gomock.Any()).Return(&Prize{}, nil)
			storageMock.EXPECT().Update(gomock.Any()).Return(mockedErr)

			res, err := manager.Edit(&PrizeEditRequest{ID: testID})

			assert.Error(t, err)
			assert.ErrorIs(t, err, mockedErr)
			assert.Equal(t, &Result{StatusError}, res)
		})
	})

	t.Run("list", func(t *testing.T) {

		t.Run("prize_success", func(t *testing.T) {
			storageMock.EXPECT().GetAll().Return([]Prize{}, nil)

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
