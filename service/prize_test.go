package service

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrize(t *testing.T) {
	ctrl := gomock.NewController(t)
	pzMock := NewMockPrizeStorage(ctrl)
	pm := NewPrizeManager(pzMock)

	prz := &PrizeRequest{
		Name:        "prize_name_1",
		TicketCost:  1234,
		Description: "prize_description_1",
	}

	mockedID := "prize_id_1"
	mockedTime := time.Now().UTC()
	mockedErr := assert.AnError

	mockedPrize := Prize{
		ID:          mockedID,
		Name:        "prize_name_1",
		TicketCost:  1234,
		Description: "prize_description_1",
		CreatedAt:   mockedTime,
	}

	t.Run("create", func(t *testing.T) {
		t.Run("error", func(t *testing.T) {
			pzMock.EXPECT().Create(gomock.Any()).Return(mockedErr)

			res, err := pm.Create(prz)
			require.ErrorIs(t, err, mockedErr)
			require.Empty(t, res)
		})

		t.Run("success", func(t *testing.T) {
			setUUIDMock(mockedID)
			setTimeNowMock(mockedTime)

			pzMock.EXPECT().Create(&mockedPrize).Return(nil)

			resID, err := pm.Create(prz)
			require.NoError(t, err)
			require.Equal(t, mockedID, resID)
		})
	})

	t.Run("get", func(t *testing.T) {
		t.Run("error", func(t *testing.T) {
			pzMock.EXPECT().Get(mockedID).Return(nil, mockedErr)

			res, err := pm.Get(mockedID)
			require.ErrorIs(t, err, mockedErr)
			require.Nil(t, res)
		})

		t.Run("success", func(t *testing.T) {
			pzMock.EXPECT().Get(mockedID).Return(&mockedPrize, nil)

			prz, err := pm.Get(mockedID)
			require.NoError(t, err)
			require.Equal(t, &mockedPrize, prz)
		})
	})

	t.Run("edit", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			pzMock.EXPECT().Get(mockedID).Return(&mockedPrize, nil)
			pzMock.EXPECT().Update(&mockedPrize).Return(nil)

			err := pm.Edit(mockedID, prz)
			require.NoError(t, err)
		})

		t.Run("not_found", func(t *testing.T) {
			pzMock.EXPECT().Get(mockedID).Return(nil, ErrNotFound)

			err := pm.Edit(mockedID, prz)
			require.ErrorIs(t, err, ErrNotFound)
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			pzMock.EXPECT().Delete(mockedID).Return(nil)

			err := pm.Delete(mockedID)
			require.NoError(t, err)
		})

		t.Run("not_found", func(t *testing.T) {
			pzMock.EXPECT().Delete(mockedID).Return(ErrNotFound)

			err := pm.Delete(mockedID)
			require.ErrorIs(t, err, ErrNotFound)
		})
	})

	t.Run("list", func(t *testing.T) {
		t.Run("error", func(t *testing.T) {
			pzMock.EXPECT().GetAll().Return(nil, mockedErr)

			res, err := pm.List()
			require.ErrorIs(t, err, mockedErr)
			require.Nil(t, res)
		})

		t.Run("success", func(t *testing.T) {
			raffles := []Prize{mockedPrize, mockedPrize, mockedPrize}

			pzMock.EXPECT().GetAll().Return(raffles, nil)

			res, err := pm.List()
			require.NoError(t, err)
			require.Equal(t, raffles, res)
		})
	})
}
