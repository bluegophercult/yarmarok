package storage

import (
	"github.com/kaznasho/yarmarok/service"

	"cloud.google.com/go/firestore"
)

const (
	organizerCollection   = "organizers"
	raffleCollection      = "raffles"
	participantCollection = "participants"
	prizeCollection       = "prizes"
	donationCollection    = "donations"
)

// FirestoreOrganizerStorage is a storage for organizers based on Firestore.
type FirestoreOrganizerStorage struct {
	*StorageBase[service.Organizer]
}

// NewFirestoreOrganizerStorage creates a new FirestoreOrganizerStorage.
func NewFirestoreOrganizerStorage(client *firestore.Client) *FirestoreOrganizerStorage {
	idExtractor := IDExtractor[service.Organizer](
		func(o *service.Organizer) string {
			return o.ID
		},
	)

	base := NewStorageBase(client.Collection(organizerCollection), idExtractor)
	return &FirestoreOrganizerStorage{
		StorageBase: base,
	}
}

// RaffleStorage returns a storage for raffles.
func (os *FirestoreOrganizerStorage) RaffleStorage(organizerID string) service.RaffleStorage {
	return NewFirestoreRaffleStorage(os.collectionReference.Doc(organizerID).Collection(raffleCollection), organizerID)
}
