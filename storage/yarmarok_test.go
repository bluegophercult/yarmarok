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

	y.UserID = u.ID
	created := []service.Yarmarok{*y}

	t.Run("get", func(t *testing.T) {
		y2, err := ys.Get(y.ID)
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

	t.Run("create another", func(t *testing.T) {
		y2 := &service.Yarmarok{
			ID:        "yarmarok_id_2",
			Name:      "yarmarok_name_2",
			Note:      "yarmarok_note_2",
			CreatedAt: time.Now().UTC().Truncate(time.Millisecond),
			UserID:    "to be replaced",
		}

		err = ys.Create(y2)
		require.NoError(t, err)

		y2.UserID = u.ID
		created = append(created, *y2)
	})

	t.Run("list", func(t *testing.T) {
		ys, err := ys.GetAll()
		require.NoError(t, err)
		require.Len(t, ys, 2)
		require.Equal(t, created, ys)
	})
}
