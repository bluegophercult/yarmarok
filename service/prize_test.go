package service_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaznasho/yarmarok/mocks"
	"github.com/kaznasho/yarmarok/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrizeManager(t *testing.T) {
	ctrl := gomock.NewController(t)

	storageMock := mocks.NewMockPrizeStorage(ctrl)
	manager := service.NewPrizeManager(storageMock)

	id := "prize_id_1"
	p := &service.PrizeRequest{
		Name:        "prize_name_1",
		TicketCost:  1234,
		Description: "prize_description_1",
	}

	t.Run("add_prize", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			storageMock.EXPECT().Create(gomock.Any()).Return(nil)

			res, err := manager.Create(p)

			require.NoError(t, err)
			require.NotNil(t, res)
		})

		t.Run("already_exists", func(t *testing.T) {
			storageMock.EXPECT().Create(gomock.Any()).Return(service.ErrAlreadyExists)

			prizeManager := service.NewPrizeManager(storageMock)

			res, err := prizeManager.Create(&service.PrizeRequest{
				Name:        "prize_name_1",
				TicketCost:  1234,
				Description: "prize_description_1",
			})

			require.Error(t, err)
			require.ErrorIs(t, err, service.ErrAlreadyExists)
			require.Nil(t, res)
		})
	})

	t.Run("edit_prize", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			storageMock.EXPECT().Get(gomock.Any()).Return(&service.Prize{}, nil)
			storageMock.EXPECT().Update(gomock.Any()).Return(nil)

			err := manager.Edit(id, p)
			require.NoError(t, err)
		})

		t.Run("not_found", func(t *testing.T) {
			storageMock.EXPECT().Get(gomock.Any()).Return(nil, service.ErrAlreadyExists)

			err := manager.Edit(id, p)
			require.ErrorIs(t, err, service.ErrAlreadyExists)
		})

		t.Run("error", func(t *testing.T) {
			mockedErr := assert.AnError

			storageMock.EXPECT().Get(gomock.Any()).Return(&service.Prize{}, nil)
			storageMock.EXPECT().Update(gomock.Any()).Return(mockedErr)

			err := manager.Edit(id, p)
			require.ErrorIs(t, err, mockedErr)
		})
	})

	t.Run("delete_prize", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			storageMock.EXPECT().Delete(gomock.Any()).Return(nil)

			err := manager.Delete(id)
			require.NoError(t, err)
		})

		t.Run("not_found", func(t *testing.T) {
			storageMock.EXPECT().Delete(gomock.Any()).Return(service.ErrNotFound)

			err := manager.Delete(id)
			require.ErrorIs(t, err, service.ErrNotFound)
		})
	})

	t.Run("list_prizes", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			storageMock.EXPECT().GetAll().Return([]service.Prize{}, nil)

			res, err := manager.List()

			require.NoError(t, err)
			require.NotNil(t, res)
		})

		t.Run("error", func(t *testing.T) {
			mockedErr := assert.AnError

			storageMock.EXPECT().GetAll().Return(nil, mockedErr)

			res, err := manager.List()

			require.Error(t, err)
			require.ErrorIs(t, err, mockedErr)
			require.Nil(t, res)
		})
	})
}
