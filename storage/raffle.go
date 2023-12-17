package storage

import (
	"cloud.google.com/go/firestore"

	"github.com/kaznasho/yarmarok/auditlog"
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
		organizerID: organizerID,
		StorageBase: NewStorageBase(client, raffleIDExtractor),
	}
}

// FirestoreRaffleStorage is a storage for raffles based on Firestore.
type FirestoreRaffleStorage struct {
	organizerID string
	*StorageBase[service.Raffle]
}

func (rs *FirestoreRaffleStorage) Create(r *service.Raffle) error {
	r.OrganizerID = rs.organizerID
	return rs.StorageBase.Create(r)
}

// PrizeStorage returns a prize storage.
func (rs *FirestoreRaffleStorage) PrizeStorage(raffleID string) service.PrizeStorage {
	return NewFirestorePrizeStorage(rs.collectionReference.Doc(raffleID).Collection(prizeCollection), raffleID)
}

// ParticipantStorage returns a participant storage.
func (rs *FirestoreRaffleStorage) ParticipantStorage(raffleID string) service.ParticipantStorage {
	return NewFirestoreParticipantStorage(rs.collectionReference.Doc(raffleID).Collection(participantCollection), raffleID)
}

func (rs *FirestoreRaffleStorage) AuditLogStorage(raffleID string) auditlog.AuditLogStorage {
	return NewFirestoreAuditLogStorage(rs.collectionReference.Doc(raffleID).Collection(auditLogCollection))
}
