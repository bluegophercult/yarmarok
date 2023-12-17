package storage

import (
	"testing"
	"time"

	"github.com/kaznasho/yarmarok/auditlog"
	"github.com/kaznasho/yarmarok/service"
	"github.com/kaznasho/yarmarok/testinfra"
	fsemulator "github.com/kaznasho/yarmarok/testinfra/firestore"
	"github.com/stretchr/testify/require"
)

func TestAuditLogStorage(t *testing.T) {
	testinfra.SkipIfNotIntegrationRun(t)

	firestoreInstance, err := fsemulator.RunInstance(t)
	require.NoError(t, err)

	orgStorage := NewFirestoreOrganizerStorage(firestoreInstance.Client())

	org := &service.Organizer{ID: "organizer_id_1"}
	err = orgStorage.Create(org)
	require.NoError(t, err)

	raffle := service.Raffle{ID: "raffle_id_1"}
	raffleStorage := NewFirestoreRaffleStorage(orgStorage.collectionReference.Doc(org.ID).Collection(raffleCollection), raffle.ID)
	err = raffleStorage.Create(&raffle)
	require.NoError(t, err)

	storage := raffleStorage.AuditLogStorage(raffle.ID)
	record := auditlog.NewAuditLogRecord("actor", "action", nil)
	record.CreatedAt = time.Now().UTC().Truncate(time.Millisecond)

	t.Run("create", func(t *testing.T) {
		err := storage.Create(record)
		require.NoError(t, err)
	})

	t.Run("already_exists", func(t *testing.T) {
		err := storage.Create(record)
		require.Error(t, err)
	})

	t.Run("list", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			records, err := storage.GetAll()
			require.NoError(t, err)

			require.Len(t, records, 1)
			require.Equal(t, *record, records[0])
		})
	})

}

var _ auditlog.AuditLogStorage = &FirestoreAuditLogStorage{}
