package storage

import (
	"testing"

	"github.com/kaznasho/yarmarok/service"
	fsemulator "github.com/kaznasho/yarmarok/testinfra/firestore"

	"github.com/kaznasho/yarmarok/testinfra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateOrganizer(t *testing.T) {
	testinfra.SkipIfNotIntegrationRun(t)
	org := service.Organizer{ID: "123"}

	firestoreInstance, err := fsemulator.RunInstance(t)
	require.NoError(t, err)

	orgStorage := NewFirestoreOrganizerStorage(firestoreInstance.Client())

	t.Run("create", func(t *testing.T) {

		err = orgStorage.Create(org)
		require.NoError(t, err)
	})

	t.Run("exists", func(t *testing.T) {
		exists, err := orgStorage.Exists(org.ID)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("not exists", func(t *testing.T) {
		exists, err := orgStorage.Exists("not-exists")
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("create again", func(t *testing.T) {
		err = orgStorage.Create(org)
		require.ErrorIs(t, err, service.ErrOrganizerAlreadyExists)
	})
}

var _ service.OrganizerStorage = &FirestoreOrganizerStorage{}
