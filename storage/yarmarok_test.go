package storage

import (
	"testing"
	"time"

	"github.com/kaznasho/yarmarok/service"
	"github.com/kaznasho/yarmarok/testinfra"
	fsemulator "github.com/kaznasho/yarmarok/testinfra/firestore"
	"github.com/stretchr/testify/require"
)

func TestYarmarok(t *testing.T) {
	testinfra.SkipIfNotIntegrationRun(t)

	firestoreInstance, err := fsemulator.RunInstance(t)
	require.NoError(t, err)

	us := NewFirestoreUserStorage(firestoreInstance.Client())

	u := service.User{ID: "user_id_1"}
	err = us.Create(u)
	require.NoError(t, err)

	ys := us.YarmarokStorage(u.ID)

	y := &service.Yarmarok{
		ID:        "yarmarok_id_1",
		Name:      "yarmarok_name_1",
		Note:      "yarmarok_note_1",
		CreatedAt: time.Now().UTC().Truncate(time.Millisecond),
		UserID:    "to be replaced",
	}

	t.Run("create", func(t *testing.T) {
		err = ys.Create(y)
		require.NoError(t, err)
	})

	t.Run("get", func(t *testing.T) {
		y2, err := ys.Get(y.ID)
		y.UserID = u.ID
		require.NoError(t, err)
		require.Equal(t, y, y2)
	})

	t.Run("not exists", func(t *testing.T) {
		resp, err := ys.Get("not-exists")
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("create again", func(t *testing.T) {
		err = ys.Create(y)
		require.ErrorIs(t, err, service.ErrYarmarokAlreadyExists)
	})

	t.Run("create without id", func(t *testing.T) {
		y2 := &service.Yarmarok{
			Name: "yarmarok_name_2",
			Note: "yarmarok_note_2",
		}
		err = ys.Create(y2)
		require.Error(t, err)
	})
}
