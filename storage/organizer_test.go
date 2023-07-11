package storage

import (
	"testing"

	"github.com/kaznasho/yarmarok/service"
	"github.com/kaznasho/yarmarok/testinfra/firestore"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaznasho/yarmarok/testinfra"
)

func TestCreateOrganizer(t *testing.T) {
	testinfra.SkipIfNotIntegrationRun(t)
	org := &service.Organizer{ID: "123"}

	firestoreInstance, err := firestore.RunInstance(t)
	require.NoError(t, err)

	os := NewFirestoreOrganizerStorage(firestoreInstance.Client())

	t.Run("create", func(t *testing.T) {
		err = os.Create(org)
		require.NoError(t, err)
	})

	t.Run("exists", func(t *testing.T) {
		exists, err := os.Exists(org.ID)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("not exists", func(t *testing.T) {
		exists, err := os.Exists("not-exists")
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("create again", func(t *testing.T) {
		err = os.Create(org)
		require.ErrorIs(t, err, service.ErrAlreadyExists)
	})
}

var _ service.OrganizerStorage = &FirestoreOrganizerStorage{}
