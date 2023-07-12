package storage

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kaznasho/yarmarok/service"
	"github.com/kaznasho/yarmarok/testinfra"
	"github.com/kaznasho/yarmarok/testinfra/firestore"
)

func TestParticipantStorage(t *testing.T) {
	testinfra.SkipIfNotIntegrationRun(t)

	firestoreInstance, err := firestore.RunInstance(t)
	require.NoError(t, err)

	os := NewFirestoreOrganizerStorage(firestoreInstance.Client())

	org := service.Organizer{ID: "organizer_id_1"}
	err = os.Create(org)
	require.NoError(t, err)

	raf := service.Raffle{ID: "raffle_id_1"}
	rs := NewFirestoreRaffleStorage(os.firestoreClient.Doc(org.ID).Collection(raffleCollection), raf.ID)

	err = rs.Create(&raf)
	require.NoError(t, err)

	ps := rs.ParticipantStorage(raf.ID)

	t.Run("Participant operations", func(t *testing.T) {
		created := make([]service.Participant, 0)

		for i := 1; i <= 5; i++ {
			p := service.Participant{
				ID:    fmt.Sprintf("participant_id_%d", i),
				Name:  fmt.Sprintf("Participant %d", i),
				Phone: fmt.Sprintf("12345678%d", i),
				Note:  fmt.Sprintf("Participant %d notes", i),
			}

			t.Run(fmt.Sprintf("Create participant %d", i), func(t *testing.T) {
				err = ps.Create(&p)
				require.NoError(t, err)
				created = append(created, p)
			})

			t.Run(fmt.Sprintf("Get participant %d", i), func(t *testing.T) {
				p2, err := ps.Get(p.ID)
				require.NoError(t, err)
				require.Equal(t, &p, p2)
			})

			t.Run(fmt.Sprintf("Update participant %d", i), func(t *testing.T) {
				p.Name = fmt.Sprintf("Updated Participant %d", i)
				err = ps.Update(&p)
				require.NoError(t, err)

				p2, err := ps.Get(p.ID)
				require.NoError(t, err)
				require.Equal(t, &p, p2)

				created[i-1] = p
			})

			t.Run("Get all participants", func(t *testing.T) {
				participants, err := ps.GetAll()
				require.NoError(t, err)
				require.ElementsMatch(t, created, participants)
			})
		}

		t.Run("Get non-existent participant", func(t *testing.T) {
			resp, err := ps.Get("not-exists")
			require.Error(t, err)
			require.Nil(t, resp)
		})
	})
}

var _ service.ParticipantStorage = (*FirestoreParticipantStorage)(nil)
