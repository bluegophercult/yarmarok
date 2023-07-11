package storage

import (
	"context"
	"fmt"

	"github.com/kaznasho/yarmarok/service"

	"cloud.google.com/go/firestore"
)

const (
	organizerCollection   = "organizers"
	raffleCollection      = "raffles"
	participantCollection = "participants"
	prizeCollection       = "prizes"
)

// FirestoreOrganizerStorage is a storage for organizers based on Firestore.
type FirestoreOrganizerStorage struct {
	firestoreClient *firestore.CollectionRef
}

// NewFirestoreOrganizerStorage creates a new FirestoreOrganizerStorage.
func NewFirestoreOrganizerStorage(client *firestore.Client) *FirestoreOrganizerStorage {
	return &FirestoreOrganizerStorage{
		firestoreClient: client.Collection(organizerCollection),
	}
}

// Create creates a new organizer.
func (os *FirestoreOrganizerStorage) Create(org service.Organizer) error {
	exists, err := os.Exists(org.ID)
	if err != nil {
		return fmt.Errorf("check organizer exists: %w", err)
	}

	if exists {
		return service.ErrOrganizerAlreadyExists
	}

	_, err = os.firestoreClient.Doc(org.ID).Set(context.Background(), org)
	return err
}

// Exists checks if an organizer with the given id exists.
func (os *FirestoreOrganizerStorage) Exists(id string) (bool, error) {
	doc, err := os.firestoreClient.Doc(id).Get(context.Background())
	if isNotFound(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return doc.Exists(), nil
}

// RaffleStorage returns a storage for raffles.
func (os *FirestoreOrganizerStorage) RaffleStorage(organizerID string) service.RaffleStorage {
	return NewFirestoreRaffleStorage(os.firestoreClient.Doc(organizerID).Collection(raffleCollection), organizerID)
}
