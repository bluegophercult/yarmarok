package service

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:generate mockgen -destination=mock_contributor_storage_test.go -package=service github.com/kaznasho/yarmarok/service ContributorStorage

func TestContributorManagerAdd(t *testing.T) {
	ctrl := gomock.NewController(t)

	storageMock := NewMockContributorStorage(ctrl)
	manager := NewContributorManager(storageMock)

	t.Run("add contributor", func(t *testing.T) {
		storageMock.EXPECT().Create(gomock.Any()).Return(nil)

		_, err := manager.Add(&ContributorAddRequest{
			Name:  "John Doe",
			Phone: "1234567890",
			Note:  "Test contributor",
		})

		assert.NoError(t, err)
	})

	t.Run("add_already_exists", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		storageMock.EXPECT().Create(gomock.Any()).Return(ErrContributorAlreadyExists)

		contributorManager := NewContributorManager(storageMock)

		_, err := contributorManager.Add(&ContributorAddRequest{
			Name:  "John Doe",
			Phone: "1234567890",
			Note:  "Test contributor",
		})

		assert.ErrorIs(t, err, ErrContributorAlreadyExists)
	})
}

func TestContributorManagerEdit(t *testing.T) {
	ctrl := gomock.NewController(t)

	storageMock := NewMockContributorStorage(ctrl)
	manager := NewContributorManager(storageMock)

	t.Run("edit contributor", func(t *testing.T) {
		storageMock.EXPECT().Get(gomock.Any()).Return(&Contributor{}, nil)
		storageMock.EXPECT().Update(gomock.Any()).Return(nil)

		_, err := manager.Edit(&ContributorEditRequest{ID: "test-id"})

		assert.NoError(t, err)
	})

	t.Run("contributor not found", func(t *testing.T) {
		storageMock.EXPECT().Get(gomock.Any()).Return(nil, ErrContributorNotFound)

		_, err := manager.Edit(&ContributorEditRequest{
			ID: "test-id",
		})

		assert.ErrorIs(t, err, ErrContributorNotFound)
	})
}

func TestContributorManagerList(t *testing.T) {
	ctrl := gomock.NewController(t)

	storageMock := NewMockContributorStorage(ctrl)
	manager := NewContributorManager(storageMock)

	t.Run("list contributors", func(t *testing.T) {
		storageMock.EXPECT().GetAll().Return([]Contributor{}, nil)

		_, err := manager.List()

		assert.NoError(t, err)
	})
}
