package storage

import (
	"fmt"
	"testing"
	"time"

	"github.com/kaznasho/yarmarok/service"
	"github.com/kaznasho/yarmarok/testinfra"
	fsemulator "github.com/kaznasho/yarmarok/testinfra/firestore"
	"github.com/stretchr/testify/require"
)

func TestDonationStorage(t *testing.T) {
	var donationCollection = "donations"

	testinfra.SkipIfNotIntegrationRun(t)

	firestoreInstance, err := fsemulator.RunInstance(t)
	require.NoError(t, err)

	participantID := "participant_id_1"
	participantStorage := NewFirestoreParticipantStorage(firestoreInstance.Client().Collection(participantCollection), participantID)
	err = participantStorage.Create(&service.Participant{ID: participantID})
	require.NoError(t, err)

	prizeID := "prize_id_1"
	prizeStorage := NewFirestorePrizeStorage(firestoreInstance.Client().Collection(prizeCollection), prizeID)
	err = prizeStorage.Create(&service.Prize{ID: prizeID})
	require.NoError(t, err)

	ds := NewFirestoreDonationStorage(firestoreInstance.Client().Doc(participantID).Collection(donationCollection), participantID)

	t.Run("Donation operations", func(t *testing.T) {
		testDonations := make([]service.Donation, 0)

		for i := 1; i <= 5; i++ {
			d := &service.Donation{
				ID:            fmt.Sprintf("donation_id_%d", i),
				PrizeID:       prizeID,
				ParticipantID: participantID,
				Amount:        100 + i,
				TicketNumber:  i,
				Description:   fmt.Sprintf("donation_description_%d", i),
				CreatedAt:     time.Now().UTC().Truncate(time.Millisecond),
			}

			t.Run("Create donation", func(t *testing.T) {
				err = ds.Create(participantStorage, prizeStorage, d)
				require.NoError(t, err)
				testDonations = append(testDonations, *d)
			})

			t.Run("Get donation", func(t *testing.T) {
				d2, err := ds.Get(d.ID)
				require.NoError(t, err)
				require.Equal(t, d, d2)
			})

			t.Run("Update donation", func(t *testing.T) {
				d.Description = fmt.Sprintf("updated_description_%d", i)
				err = ds.Update(d)
				require.NoError(t, err)

				d2, err := ds.Get(d.ID)
				require.NoError(t, err)
				require.Equal(t, d, d2)

				testDonations[i-1] = *d
			})

			t.Run("Get all donations", func(t *testing.T) {
				getDonations, err := ds.GetAll()
				require.NoError(t, err)
				require.ElementsMatch(t, testDonations, getDonations)
			})
		}

		t.Run("Get non-existent donation", func(t *testing.T) {
			resp, err := ds.Get("not-exists")
			require.Error(t, err)
			require.Nil(t, resp)
		})
	})
}
