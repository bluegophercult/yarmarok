package storage

import (
	"cloud.google.com/go/firestore"

	"github.com/kaznasho/yarmarok/service"
)

// NewFirestoreRaffleStorage creates a new FirestoreRaffleStorage.
func NewFirestoreRaffleStorage(client *firestore.CollectionRef, organizerID string) *FirestoreRaffleStorage {
	raffleIDExtractor := IDExtractor[service.Raffle](
		func(r *service.Raffle) string {
			return r.ID
		},
	)

	return &FirestoreRaffleStorage{
		StorageBase: NewStorageBase(client, raffleIDExtractor),
	}
}

// FirestoreRaffleStorage is a storage for raffles based on Firestore.
type FirestoreRaffleStorage struct {
	*StorageBase[service.Raffle]
}

// PrizeStorage returns a prize storage.
func (rs *FirestoreRaffleStorage) PrizeStorage(raffleID string) service.PrizeStorage {
	return NewFirestorePrizeStorage(rs.collectionReference.Doc(raffleID).Collection(prizeCollection), raffleID)
}

// ParticipantStorage returns a participant storage.
func (rs *FirestoreRaffleStorage) ParticipantStorage(raffleID string) service.ParticipantStorage {
	return NewFirestoreParticipantStorage(rs.collectionReference.Doc(raffleID).Collection(participantCollection), raffleID)
}
