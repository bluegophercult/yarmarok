package storage

import (
	"testing"

	"github.com/kaznasho/yarmarok/service"
	fsemulator "github.com/kaznasho/yarmarok/testinfra/firestore"

	"github.com/kaznasho/yarmarok/testinfra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	testinfra.SkipIfNotIntegrationRun(t)
	u := service.User{ID: "123"}

	firestoreInstance, err := fsemulator.RunInstance(t)
	require.NoError(t, err)

	us := NewFirestoreUserStorage(firestoreInstance.Client())

	t.Run("create", func(t *testing.T) {

		err = us.Create(u)
		require.NoError(t, err)
	})

	t.Run("exists", func(t *testing.T) {
		exists, err := us.Exists(u.ID)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("not exists", func(t *testing.T) {
		exists, err := us.Exists("not-exists")
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("create again", func(t *testing.T) {
		err = us.Create(u)
		require.ErrorIs(t, err, service.ErrUserAlreadyExists)
	})
}

var _ service.UserStorage = &FirestoreUserStorage{}
