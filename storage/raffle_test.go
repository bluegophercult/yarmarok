package storage

import (
	"testing"
	"time"

	"github.com/kaznasho/yarmarok/service"
	"github.com/kaznasho/yarmarok/testinfra"
	fsemulator "github.com/kaznasho/yarmarok/testinfra/firestore"
	"github.com/stretchr/testify/require"
)

func TestRaffle(t *testing.T) {
	testinfra.SkipIfNotIntegrationRun(t)

	firestoreInstance, err := fsemulator.RunInstance(t)
	require.NoError(t, err)

	orgStorage := NewFirestoreOrganizerStorage(firestoreInstance.Client())

	org := service.Organizer{ID: "organizer_id_1"}
	err = orgStorage.Create(org)
	require.NoError(t, err)

	rafStorage := orgStorage.RaffleStorage(org.ID)

	raf := &service.Raffle{
		ID:          "raffle_id_1",
		Name:        "raffle_name_1",
		Note:        "raffle_note_1",
		CreatedAt:   time.Now().UTC().Truncate(time.Millisecond),
		OrganizerID: "to be replaced",
	}

	t.Run("create", func(t *testing.T) {
		err = rafStorage.Create(raf)
		require.NoError(t, err)
	})

	raf.OrganizerID = org.ID
	created := []service.Raffle{*raf}

	t.Run("get", func(t *testing.T) {
		raf2, err := rafStorage.Get(raf.ID)
		require.NoError(t, err)
		require.Equal(t, raf, raf2)
	})

	t.Run("not exists", func(t *testing.T) {
		resp, err := rafStorage.Get("not-exists")
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("create again", func(t *testing.T) {
		err = rafStorage.Create(raf)
		require.ErrorIs(t, err, service.ErrRaffleAlreadyExists)
	})

	t.Run("create without id", func(t *testing.T) {
		raf2 := &service.Raffle{
			Name: "raffle_name_2",
			Note: "raffle_note_2",
		}
		err = rafStorage.Create(raf2)
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

		err = rafStorage.Create(raf2)
		require.NoError(t, err)

		raf2.OrganizerID = org.ID
		created = append(created, *raf2)
	})

	t.Run("list", func(t *testing.T) {
		raffles, err := rafStorage.GetAll()
		require.NoError(t, err)
		require.Len(t, raffles, 2)
		require.Equal(t, created, raffles)
	})
}
