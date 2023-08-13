package storage

import (
	"testing"
	"time"

	"github.com/kaznasho/yarmarok/testinfra"
	"github.com/stretchr/testify/require"

	"github.com/kaznasho/yarmarok/service"
	"github.com/kaznasho/yarmarok/testinfra/firestore"
)

func TestRaffle(t *testing.T) {
	testinfra.SkipIfNotIntegrationRun(t)

	firestoreInstance, err := firestore.RunInstance(t)
	require.NoError(t, err)

	os := NewFirestoreOrganizerStorage(firestoreInstance.Client())

	org := &service.Organizer{ID: "organizer_id_1"}
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

	t.Run("query", func(t *testing.T) {
		raf3 := &service.Raffle{ID: "raffle_id_3", Name: "E Raffle", Note: "note3", CreatedAt: time.Now().UTC().Truncate(time.Millisecond)}
		raf4 := &service.Raffle{ID: "raffle_id_4", Name: "F Raffle", Note: "note4", CreatedAt: time.Now().UTC().Truncate(time.Millisecond)}
		raf5 := &service.Raffle{ID: "raffle_id_5", Name: "G Raffle", Note: "note5", CreatedAt: time.Now().UTC().Truncate(time.Millisecond)}
		raf6 := &service.Raffle{ID: "raffle_id_6", Name: "H Raffle", Note: "note6", CreatedAt: time.Now().UTC().Truncate(time.Millisecond)}

		err = rs.Create(raf3)
		require.NoError(t, err)
		err = rs.Create(raf4)
		require.NoError(t, err)
		err = rs.Create(raf5)
		require.NoError(t, err)
		err = rs.Create(raf6)
		require.NoError(t, err)

		query := new(service.Query).
			WithFilter("Name", service.GT, "C Raffle").
			WithOrderBy("Name", service.ASC).
			WithLimit(3)

		raffles, err := rs.Query(query)
		require.NoError(t, err)
		require.Equal(t, []service.Raffle{*raf3, *raf4, *raf5}, raffles)

		query = new(service.Query).
			WithFilter("Name", service.GTE, "E Raffle").
			WithFilter("Note", service.EQ, "note3").
			WithOrderBy("Name", service.DESC).
			WithLimit(2)

		raffles, err = rs.Query(query)
		require.NoError(t, err)
		require.Equal(t, []service.Raffle{*raf3}, raffles)

		query = new(service.Query).WithFilter("ID", service.IN, []string{"fake_id_1"})
		raffles, err = rs.Query(query)
		require.ErrorIs(t, err, service.ErrNotFound)
	})
}
