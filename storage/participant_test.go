package storage

import (
	"fmt"
	"github.com/kaznasho/yarmarok/testinfra"
	"testing"

	"github.com/kaznasho/yarmarok/service"
	fsemulator "github.com/kaznasho/yarmarok/testinfra/firestore"
	"github.com/stretchr/testify/require"
)

func TestParticipantStorage(t *testing.T) {
	testinfra.SkipIfNotIntegrationRun(t)

	firestoreInstance, err := fsemulator.RunInstance(t)
	require.NoError(t, err)

	us := NewFirestoreUserStorage(firestoreInstance.Client())

	u := service.User{ID: "user_id_1"}
	err = us.Create(u)
	require.NoError(t, err)

	ps := us.ParticipantStorage(u.ID)

	p := service.Participant{
		ID:         "participant_id_1",
		YarmarokID: "yarmarok_id_1",
		Name:       "Participant 1",
		Phone:      "123456789",
		Email:      "participant1@example.com",
		Notes:      "Participant 1 notes",
	}

	t.Run("create", func(t *testing.T) {
		err = ps.Create(p)
		require.NoError(t, err)
	})

	created := []service.Participant{p}

	t.Run("get", func(t *testing.T) {
		p2, err := ps.Get(p.ID)
		require.NoError(t, err)
		require.Equal(t, &p, p2)
	})

	t.Run("not exists", func(t *testing.T) {
		resp, err := ps.Get("not-exists")
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("create again", func(t *testing.T) {
		err = ps.Create(p)
		require.ErrorIs(t, err, service.ErrParticipantAlreadyExists)
	})

	t.Run("update", func(t *testing.T) {
		p.Name = "Updated Participant 1"
		err = ps.Update(p)
		require.NoError(t, err)

		p2, err := ps.Get(p.ID)
		require.NoError(t, err)
		require.Equal(t, &p, p2)
	})

	t.Run("delete", func(t *testing.T) {
		err = ps.Delete(p.ID)
		require.NoError(t, err)

		resp, err := ps.Get(p.ID)
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("get all", func(t *testing.T) {
		for i := 2; i <= 5; i++ {
			p := service.Participant{
				ID:         fmt.Sprintf("participant_id_%d", i),
				YarmarokID: "yarmarok_id_1",
				Name:       fmt.Sprintf("Participant %d", i),
				Phone:      fmt.Sprintf("12345678%d", i),
				Email:      fmt.Sprintf("participant%d@example.com", i),
				Notes:      fmt.Sprintf("Participant %d notes", i),
			}
			err = ps.Create(p)
			require.NoError(t, err)
			created = append(created, p)
		}

		participants, err := ps.GetAll()
		require.NoError(t, err)
		require.ElementsMatch(t, created, participants)
	})
}

var _ service.ParticipantStorage = (*FirestoreParticipantStorage)(nil)
