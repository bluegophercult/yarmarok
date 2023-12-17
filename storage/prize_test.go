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

func TestPrizeStorage(t *testing.T) {
	testinfra.SkipIfNotIntegrationRun(t)

	firestoreInstance, err := fsemulator.RunInstance(t)
	require.NoError(t, err)

	os := NewFirestoreOrganizerStorage(firestoreInstance.Client())

	org := &service.Organizer{ID: "organizer_id_1"}
	err = os.Create(org)
	require.NoError(t, err)

	y := service.Raffle{ID: "raffle_id_1"}
	ys := NewFirestoreRaffleStorage(os.collectionReference.Doc(org.ID).Collection(raffleCollection), y.ID)

	err = ys.Create(&y)
	require.NoError(t, err)

	pz := ys.PrizeStorage(y.ID)

	t.Run("Prize operations", func(t *testing.T) {
		testPrizes := make([]service.Prize, 0)

		for i := 1; i <= 5; i++ {
			p := service.Prize{
				ID:          fmt.Sprintf("prize_id_%d", i),
				Name:        fmt.Sprintf("prize_name_%d", i),
				TicketCost:  123 + i,
				Description: fmt.Sprintf("prize_description_%d", i),
				CreatedAt:   time.Now().UTC().Truncate(time.Millisecond),
			}

			t.Run("Create prize", func(t *testing.T) {
				err = pz.Create(&p)
				require.NoError(t, err)
				testPrizes = append(testPrizes, p)
			})

			t.Run("Get prize", func(t *testing.T) {
				p2, err := pz.Get(p.ID)
				require.NoError(t, err)
				require.Equal(t, &p, p2)
			})

			t.Run("Update prize", func(t *testing.T) {
				p.Name = fmt.Sprintf("updated_prize %d", i)
				err = pz.Update(&p)
				require.NoError(t, err)

				p2, err := pz.Get(p.ID)
				require.NoError(t, err)
				require.Equal(t, &p, p2)

				testPrizes[i-1] = p
			})

			t.Run("Get all prizes", func(t *testing.T) {
				getPrizes, err := pz.GetAll()
				require.NoError(t, err)
				require.ElementsMatch(t, testPrizes, getPrizes)
			})
		}

		t.Run("Get non-existent prize", func(t *testing.T) {
			resp, err := pz.Get("not-exists")
			require.Error(t, err)
			require.Nil(t, resp)
		})
	})
}

var _ service.PrizeStorage = &FirestorePrizeStorage{}
