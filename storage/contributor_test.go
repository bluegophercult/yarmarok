package storage

import (
	"fmt"
	"github.com/kaznasho/yarmarok/service"
	"github.com/kaznasho/yarmarok/testinfra"
	fsemulator "github.com/kaznasho/yarmarok/testinfra/firestore"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestContributorStorage(t *testing.T) {
	testinfra.SkipIfNotIntegrationRun(t)

	firestoreInstance, err := fsemulator.RunInstance(t)
	require.NoError(t, err)

	orgStorage := NewFirestoreOrganizerStorage(firestoreInstance.Client())
	org := service.Organizer{ID: "organizer_id_1"}
	err = orgStorage.Create(org)
	require.NoError(t, err)

	raf := service.Raffle{ID: "raffle_id_1"}
	rafStorage := NewFirestoreRaffleStorage(orgStorage.firestoreClient.Doc(org.ID).Collection(raffleCollection), raf.ID)

	err = rafStorage.Create(&raf)
	require.NoError(t, err)

	storage := rafStorage.ContributorStorage(raf.ID)

	t.Run("Contributor operations", func(t *testing.T) {
		created := make([]service.Contributor, 0)

		for i := 1; i <= 5; i++ {
			ctb := service.Contributor{
				ID:    fmt.Sprintf("contributor_id_%d", i),
				Name:  fmt.Sprintf("Contributor %d", i),
				Phone: fmt.Sprintf("12345678%d", i),
				Note:  fmt.Sprintf("Contributor %d notes", i),
			}

			t.Run(fmt.Sprintf("Create contributor %d", i), func(t *testing.T) {
				err = storage.Create(&ctb)
				require.NoError(t, err)
				created = append(created, ctb)
			})

			t.Run(fmt.Sprintf("Get contributor %d", i), func(t *testing.T) {
				ctb2, err := storage.Get(ctb.ID)
				require.NoError(t, err)
				require.Equal(t, &ctb, ctb2)
			})

			t.Run(fmt.Sprintf("Update contributor %d", i), func(t *testing.T) {
				ctb.Name = fmt.Sprintf("Updated Contributor %d", i)
				err = storage.Update(&ctb)
				require.NoError(t, err)

				ctb2, err := storage.Get(ctb.ID)
				require.NoError(t, err)
				require.Equal(t, &ctb, ctb2)

				created[i-1] = ctb
			})

			t.Run("Get all contributors", func(t *testing.T) {
				contributors, err := storage.GetAll()
				require.NoError(t, err)
				require.ElementsMatch(t, created, contributors)
			})
		}

		t.Run("Get non-existent contributor", func(t *testing.T) {
			resp, err := storage.Get("not-exists")
			require.Error(t, err)
			require.Nil(t, resp)
		})
	})
}

var _ service.ContributorStorage = (*FirestoreContributorStorage)(nil)
