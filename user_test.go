package yarmarok

import (
	"testing"

	fsemulator "github.com/kaznasho/yarmarok/testinfra/firestore"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	//testinfra.SkipIfNotIntegrationRun(t)
	user := User{ID: "123"}

	firestoreInstance, err := fsemulator.RunInstance(t)
	require.NoError(t, err)

	us := NewFirestoreUserStorage(firestoreInstance.Client())

	t.Run("create", func(t *testing.T) {

		err = us.Create(user)
		require.NoError(t, err)
	})

	t.Run("exists", func(t *testing.T) {
		exists, err := us.Exists(user.ID)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("not exists", func(t *testing.T) {
		exists, err := us.Exists("not-exists")
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("create again", func(t *testing.T) {
		err = us.Create(user)
		require.ErrorIs(t, err, ErrUserAlreadyExists)
	})
}

var _ UserStorage = &FirestoreUserStorage{}
