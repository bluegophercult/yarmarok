package storage

import (
	"cloud.google.com/go/firestore"

	"github.com/kaznasho/yarmarok/service"
)

// FirestorePrizeStorage is a storage for prizes based on Firestore.
type FirestorePrizeStorage struct {
	yarmarokID      string
	firestoreClient *firestore.CollectionRef
}

// PrizeStorage returns a prize storage.
func (ys *FirestoreRaffleStorage) PrizeStorage(yarmarokID string) service.PrizeStorage {
	return NewFirestorePrizeStorage(ys.firestoreClient.Doc(yarmarokID).Collection(prizeCollection), yarmarokID)
}

// NewFirestorePrizeStorage creates a new FirestorePrizeStorage.
func NewFirestorePrizeStorage(client *firestore.CollectionRef, raffleID string) *FirestorePrizeStorage {
	prizeIDExtractor := IDExtractor[service.Prize](
		func(p *service.Prize) string {
			return p.ID
		},
	)

	return &FirestorePrizeStorage{
		raffleID:    raffleID,
		StorageBase: NewStorageBase(client, prizeIDExtractor),
	}
}
