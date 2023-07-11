package storage

import (
	"cloud.google.com/go/firestore"

	"github.com/kaznasho/yarmarok/service"
)

// FirestorePrizeStorage is a storage for prizes based on Firestore.
type FirestorePrizeStorage struct {
	raffleID string
	*StorageBase[service.Prize]
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
