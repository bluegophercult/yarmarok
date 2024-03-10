package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/assert"
)

func TestDonationManagerCreateDonation(t *testing.T) {
	ctrl := gomock.NewController(t)
	storageMock := NewMockDonationStorage(ctrl)

	manager := NewDonationManager(storageMock)

	t.Run("Add donation", func(t *testing.T) {
		storageMock.EXPECT().Create(gomock.Any()).Return(nil)

		_, err := manager.Create(&DonationRequest{Amount: 777, ParticipantID: stringUUID()})
		require.NoError(t, err)
	})

	t.Run("Create error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		storageMock.EXPECT().Create(gomock.Any()).Return(assert.AnError)

		donationManager := NewDonationManager(storageMock)
		_, err := donationManager.Create(&DonationRequest{Amount: 777, ParticipantID: stringUUID()})
		require.ErrorIs(t, err, assert.AnError)
	})
}

func TestDonationManagerEditDonation(t *testing.T) {
	ctrl := gomock.NewController(t)
	storageMock := NewMockDonationStorage(ctrl)

	manager := NewDonationManager(storageMock)
	testID := "donation_test_id"

	t.Run("Edit donation", func(t *testing.T) {
		donationRequest := &DonationRequest{Amount: 999, ParticipantID: "participant_test_id"}
		donation := &Donation{ParticipantID: "participant_test_id", Amount: 999}
		storageMock.EXPECT().Get(testID).Return(&Donation{}, nil)
		storageMock.EXPECT().Update(donation).Return(nil)

		err := manager.Edit(testID, donationRequest)
		require.NoError(t, err)
	})

	t.Run("Edit not found donation", func(t *testing.T) {
		storageMock.EXPECT().Get(testID).Return(nil, ErrNotFound)

		err := manager.Edit(testID, &DonationRequest{Amount: 999, ParticipantID: "participant_test_id"})
		require.ErrorIs(t, err, ErrNotFound)
	})
}

func TestDonationManagerListDonations(t *testing.T) {
	ctrl := gomock.NewController(t)
	storageMock := NewMockDonationStorage(ctrl)
	manager := NewDonationManager(storageMock)

	t.Run("Success", func(t *testing.T) {
		date := time.Now()
		donations := []Donation{
			{ID: "1", ParticipantID: "1", Amount: 10, CreatedAt: date},
			{ID: "2", ParticipantID: "2", Amount: 20, CreatedAt: date.Add(time.Second)},
		}
		storageMock.EXPECT().GetAll().Return(donations, nil)

		res, err := manager.List()
		require.NoError(t, err)
		require.Equal(t, donations, res)
	})

	t.Run("Error", func(t *testing.T) {
		storageMock.EXPECT().GetAll().Return(nil, assert.AnError)

		res, err := manager.List()
		require.ErrorIs(t, err, assert.AnError)
		require.Nil(t, res)
	})
}

func TestDonationManagerGetDonations(t *testing.T) {
	ctrl := gomock.NewController(t)
	storageMock := NewMockDonationStorage(ctrl)
	manager := NewDonationManager(storageMock)

	t.Run("Success", func(t *testing.T) {
		donation := &Donation{ID: "1", ParticipantID: "1", Amount: 10, CreatedAt: time.Now()}
		storageMock.EXPECT().Get(donation.ID).Return(donation, nil)

		res, err := manager.Get(donation.ID)
		require.NoError(t, err)
		require.Equal(t, donation, res)
	})

	t.Run("Error", func(t *testing.T) {
		id := "donation_id"
		storageMock.EXPECT().Get(id).Return(nil, ErrNotFound)

		res, err := manager.Get(id)
		require.ErrorIs(t, err, ErrNotFound)
		require.Nil(t, res)
	})
}

func TestDonationManagerDeleteDonation(t *testing.T) {
	ctrl := gomock.NewController(t)
	storageMock := NewMockDonationStorage(ctrl)
	manager := NewDonationManager(storageMock)

	t.Run("Success", func(t *testing.T) {
		id := "donation_id"
		storageMock.EXPECT().Delete(id).Return(nil)

		err := manager.Delete(id)
		require.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		id := "donation_id"
		storageMock.EXPECT().Delete(id).Return(ErrNotFound)

		err := manager.Delete(id)
		require.ErrorIs(t, err, ErrNotFound)
	})
}
