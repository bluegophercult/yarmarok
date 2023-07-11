package storage

import (
	"cloud.google.com/go/firestore"

	"github.com/kaznasho/yarmarok/service"
)

// NewFirestoreRaffleStorage creates a new FirestoreRaffleStorage.
func NewFirestoreRaffleStorage(client *firestore.CollectionRef, organizerID string) *FirestoreRaffleStorage {
	return &FirestoreRaffleStorage{
		organizerID:     organizerID,
		firestoreClient: client,
	}
}

// FirestoreRaffleStorage is a storage for raffles based on Firestore.
type FirestoreRaffleStorage struct {
	organizerID     string
	firestoreClient *firestore.CollectionRef
	StorageBase[service.Raffle]
}
