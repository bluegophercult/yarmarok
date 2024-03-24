package storage

import (
	"cloud.google.com/go/firestore"

	"github.com/kaznasho/yarmarok/service"
)

// FirestoreDonationStorage is a storage for donation based on Firestore.
type FirestoreDonationStorage struct {
	*StorageBase[service.Donation]
}

// NewFirestoreDonationStorage creates a new FirestoreDonationStorage.
func NewFirestoreDonationStorage(client *firestore.CollectionRef) *FirestoreDonationStorage {
	donationIDExtractor := IDExtractor[service.Donation](
		func(p *service.Donation) string {
			return p.ID
		},
	)

	return &FirestoreDonationStorage{
		StorageBase: NewStorageBase(client, donationIDExtractor),
	}
}
