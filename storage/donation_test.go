package storage

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/kaznasho/yarmarok/service"
	"github.com/kaznasho/yarmarok/testinfra"
	fsemulator "github.com/kaznasho/yarmarok/testinfra/firestore"
)

func TestDonationStorage(t *testing.T) {
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

	prizeStorage := raffleStorage.PrizeStorage(raffle.ID)

	prize := service.Prize{ID: "prize_id_1", TicketCost: 10}
	err = prizeStorage.Create(&prize)
	require.NoError(t, err)

	donationStorage := prizeStorage.DonationStorage(prize.ID)

	t.Run("Donation operations", func(t *testing.T) {
		testDonations := make([]service.Donation, 0)

		for i := 1; i <= 5; i++ {
			d := &service.Donation{
				ID:            fmt.Sprintf("donation_id_%d", i),
				PrizeID:       prize.ID,
				ParticipantID: fmt.Sprintf("participant_%d", i),
				Amount:        100 + i,
				TicketsNumber: (100 + i) / prize.TicketCost,
				CreatedAt:     time.Now().UTC().Truncate(time.Millisecond),
			}

			t.Run("Create donation", func(t *testing.T) {
				err = donationStorage.Create(d)
				require.NoError(t, err)
				testDonations = append(testDonations, *d)
			})

			t.Run("Get donation", func(t *testing.T) {
				d2, err := donationStorage.Get(d.ID)
				require.NoError(t, err)
				require.Equal(t, d, d2)
			})

			t.Run("Update donation", func(t *testing.T) {
				err = donationStorage.Update(d)
				require.NoError(t, err)

				d2, err := donationStorage.Get(d.ID)
				require.NoError(t, err)
				require.Equal(t, d, d2)

				testDonations[i-1] = *d
			})

			t.Run("Get all donations", func(t *testing.T) {
				getDonations, err := donationStorage.GetAll()
				require.NoError(t, err)
				require.ElementsMatch(t, testDonations, getDonations)
			})
		}

		t.Run("Get non-existent donation", func(t *testing.T) {
			resp, err := donationStorage.Get("not-exists")
			require.Error(t, err)
			require.Nil(t, resp)
		})
	})
}
