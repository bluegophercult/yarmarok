package service

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
)

//go:generate mockgen -destination=mock_donation_storage_test.go -package=service github.com/kaznasho/yarmarok/service DonationStorage

func TestDonationManagerAddDonation(t *testing.T) {
	ctrl := gomock.NewController(t)

	storageMock := NewMockDonationStorage(ctrl)
	participantStorageMock := NewMockParticipantStorage(ctrl)
	prizeStorageMock := NewMockPrizeStorage(ctrl)

	manager := NewDonationManager(storageMock, participantStorageMock, prizeStorageMock)

	testAmount := 100
	testDescription := "Test description in donation"

	t.Run("add donation", func(t *testing.T) {
		storageMock.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		_, err := manager.AddDonation(&DonationAddRequest{
			Amount:      testAmount,
			Description: testDescription,
		})

		assert.NoError(t, err)
	})
	t.Run("add_already_exists", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		storageMock.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(ErrDonationAlreadyExists)

		donationManager := NewDonationManager(storageMock, participantStorageMock, prizeStorageMock)

		_, err := donationManager.AddDonation(&DonationAddRequest{
			Amount:      testAmount,
			Description: testDescription,
		})

		assert.ErrorIs(t, err, ErrDonationAlreadyExists)
	})
}

func TestDonationManagerEditDonation(t *testing.T) {
	ctrl := gomock.NewController(t)

	storageMock := NewMockDonationStorage(ctrl)
	participantStorageMock := NewMockParticipantStorage(ctrl)
	prizeStorageMock := NewMockPrizeStorage(ctrl)

	manager := NewDonationManager(storageMock, participantStorageMock, prizeStorageMock)

	testID := "donation_test_id"

	t.Run("edit donation", func(t *testing.T) {
		storageMock.EXPECT().Get(gomock.Any()).Return(&Donation{}, nil)
		storageMock.EXPECT().Update(gomock.Any()).Return(nil)

		_, err := manager.EditDonation(&DonationEditRequest{ID: testID})

		assert.NoError(t, err)
	})

	t.Run("donation not found", func(t *testing.T) {
		storageMock.EXPECT().Get(gomock.Any()).Return(nil, ErrDonationNotFound)

		_, err := manager.EditDonation(&DonationEditRequest{ID: testID})

		assert.ErrorIs(t, err, ErrDonationNotFound)
	})
}

func TestDonationManagerListDonation(t *testing.T) {
	ctrl := gomock.NewController(t)

	storageMock := NewMockDonationStorage(ctrl)
	participantStorageMock := NewMockParticipantStorage(ctrl)
	prizeStorageMock := NewMockPrizeStorage(ctrl)

	manager := NewDonationManager(storageMock, participantStorageMock, prizeStorageMock)

	t.Run("list donations", func(t *testing.T) {
		storageMock.EXPECT().GetAll().Return([]Donation{}, nil)

		_, err := manager.ListDonation()

		assert.NoError(t, err)
	})
}
