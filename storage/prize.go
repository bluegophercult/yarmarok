package storage

import (
	"cloud.google.com/go/firestore"

	"github.com/kaznasho/yarmarok/service"
)

// FirestorePrizeStorage is a storage for prizes based on Firestore.
type FirestorePrizeStorage struct {
	*StorageBase[service.Prize]
}

// NewFirestorePrizeStorage creates a new FirestorePrizeStorage.
func NewFirestorePrizeStorage(client *firestore.CollectionRef) *FirestorePrizeStorage {
	prizeIDExtractor := IDExtractor[service.Prize](
		func(p *service.Prize) string {
			return p.ID
		},
	)

	return &FirestorePrizeStorage{
		StorageBase: NewStorageBase(client, prizeIDExtractor),
	}
}

// DonationStorage returns a donation storage.
func (ps *FirestorePrizeStorage) DonationStorage(prizeID string) service.DonationStorage {
	return NewFirestoreDonationStorage(ps.collectionReference.Doc(prizeID).Collection(donationCollection))
}
