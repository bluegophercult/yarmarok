package storage

import (
	"cloud.google.com/go/firestore"

	"github.com/kaznasho/yarmarok/service"
)

// FirestoreDonationStorage is a storage for donation based on Firestore.
type FirestoreDonationStorage struct {
	prizeID      string
	prizeStorage service.PrizeStorage
	*StorageBase[service.Donation]
}

// NewFirestoreDonationStorage creates a new FirestoreDonationStorage.
func NewFirestoreDonationStorage(client *firestore.CollectionRef, prizeStorage service.PrizeStorage, prizeID string) *FirestoreDonationStorage {
	donationIDExtractor := IDExtractor[service.Donation](
		func(p *service.Donation) string {
			return p.ID
		},
	)

	return &FirestoreDonationStorage{
		prizeID:      prizeID,
		prizeStorage: prizeStorage,
		StorageBase:  NewStorageBase(client, donationIDExtractor),
	}
}

// Create creates a new donation in the underlying storage.
func (ds *FirestoreDonationStorage) Create(d *service.Donation) error {
	prize, err := ds.prizeStorage.Get(ds.prizeID)
	if err != nil {
		return err
	}

	d.PrizeID = prize.ID
	d.TicketsNumber = d.Amount / prize.TicketCost

	return ds.StorageBase.Create(d)
}
