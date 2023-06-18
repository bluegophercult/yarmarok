package storage

import (
	"context"
	"fmt"

	"github.com/kaznasho/yarmarok/service"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	organizerCollection   = "organizers"
	raffleCollection      = "raffles"
	contributorCollection = "contributors"
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
func (s *FirestoreOrganizerStorage) Create(org service.Organizer) error {
	exists, err := s.Exists(org.ID)
	if err != nil {
		return fmt.Errorf("check organizer exists: %w", err)
	}

	if exists {
		return service.ErrOrganizerAlreadyExists
	}

	_, err = s.firestoreClient.Doc(org.ID).Set(context.Background(), org)
	return err
}

// Exists checks if an organizer with the given id exists.
func (s *FirestoreOrganizerStorage) Exists(id string) (bool, error) {
	doc, err := s.firestoreClient.Doc(id).Get(context.Background())
	if isNotFound(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return doc.Exists(), nil
}

func isNotFound(err error) bool {
	return status.Code(err) == codes.NotFound
}
