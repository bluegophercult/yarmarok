package storage

import (
	"cloud.google.com/go/firestore"

	"github.com/kaznasho/yarmarok/service"
)

// FirestoreParticipantStorage is a storage for raffles based on Firestore.
type FirestoreParticipantStorage struct {
	*StorageBase[service.Participant]
}

// NewFirestoreParticipantStorage creates a new FirestoreParticipantStorage.
func NewFirestoreParticipantStorage(client *firestore.CollectionRef) *FirestoreParticipantStorage {
	participantIDExtractor := IDExtractor[service.Participant](
		func(p *service.Participant) string {
			return p.ID
		},
	)

	return &FirestoreParticipantStorage{
		StorageBase: NewStorageBase(client, participantIDExtractor),
	}
}
