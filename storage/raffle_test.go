package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/kaznasho/yarmarok/service"
	"github.com/kaznasho/yarmarok/testinfra"
	"github.com/kaznasho/yarmarok/testinfra/firestore"
)

func TestRaffle(t *testing.T) {
	testinfra.SkipIfNotIntegrationRun(t)

	firestoreInstance, err := firestore.RunInstance(t)
	require.NoError(t, err)

	os := NewFirestoreOrganizerStorage(firestoreInstance.Client())

	org := service.Organizer{ID: "organizer_id_1"}
	err = os.Create(org)
	require.NoError(t, err)

	rs := os.RaffleStorage(org.ID)

	raf := &service.Raffle{
		ID:          "raffle_id_1",
		OrganizerID: "to be replaced",
		Name:        "raffle_name_1",
		Note:        "raffle_note_1",
		CreatedAt:   time.Now().UTC().Truncate(time.Millisecond),
	}

	t.Run("create", func(t *testing.T) {
		err = rs.Create(raf)
		require.NoError(t, err)
	})

	raf.OrganizerID = org.ID
	created := []service.Raffle{*raf}

	t.Run("get", func(t *testing.T) {
		raf2, err := rs.Get(raf.ID)
		require.NoError(t, err)
		require.Equal(t, raf, raf2)
	})

	t.Run("not exists", func(t *testing.T) {
		resp, err := rs.Get("not-exists")
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("create again", func(t *testing.T) {
		err = rs.Create(raf)
		require.ErrorIs(t, err, service.ErrAlreadyExists)
	})

	t.Run("create without id", func(t *testing.T) {
		raf2 := &service.Raffle{
			Name: "raffle_name_2",
			Note: "raffle_note_2",
		}
		err = rs.Create(raf2)
		require.Error(t, err)
	})

	t.Run("create another", func(t *testing.T) {
		raf2 := &service.Raffle{
			ID:          "raffle_id_2",
			Name:        "raffle_name_2",
			Note:        "raffle_note_2",
			CreatedAt:   time.Now().UTC().Truncate(time.Millisecond),
			OrganizerID: "to be replaced",
		}

		err = rs.Create(raf2)
		require.NoError(t, err)

		raf2.OrganizerID = org.ID
		created = append(created, *raf2)
	})

	t.Run("list", func(t *testing.T) {
		raffles, err := rs.GetAll()
		require.NoError(t, err)
		require.Len(t, raffles, 2)
		require.Equal(t, created, raffles)
	})
}
