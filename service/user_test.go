package service

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

//go:generate mockgen -destination=mock_user_test.go -package=service github.com/kaznasho/yarmarok/service UserStorage
func TestInitUser(t *testing.T) {
	ctrl := gomock.NewController(t)

	us := NewMockUserStorage(ctrl)

	t.Run("init user", func(t *testing.T) {
		t.Run("exists", func(t *testing.T) {
			userID := "123"
			us.EXPECT().Exists(userID).Return(true, nil)
			um := NewUserManager(us)

			err := um.InitUserIfNotExists(userID)
			assert.NoError(t, err)
		})

		t.Run("not exists", func(t *testing.T) {
			userID := "123"
			us.EXPECT().Exists(userID).Return(false, nil)
			us.EXPECT().Create(User{ID: userID}).Return(nil)
			um := NewUserManager(us)

			err := um.InitUserIfNotExists(userID)
			assert.NoError(t, err)
		})

		t.Run("error", func(t *testing.T) {
			userID := "123"
			us.EXPECT().Exists(userID).Return(false, assert.AnError)
			um := NewUserManager(us)

			err := um.InitUserIfNotExists(userID)
			assert.Error(t, err)
		})
	})
}
