package storage

import (
	"cloud.google.com/go/firestore"

	"github.com/kaznasho/yarmarok/service"
)

type FirestoreDonationStorage struct {
	prizeID      string
	prizeStorage service.PrizeStorage
	*StorageBase[service.Donation]
}

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

func (ds *FirestoreDonationStorage) Create(d *service.Donation) error {
	prize, err := ds.prizeStorage.Get(ds.prizeID)
	if err != nil {
		return err
	}

	d.PrizeID = prize.ID
	d.TicketsNumber = d.Amount / prize.TicketCost

	return ds.StorageBase.Create(d)
}
