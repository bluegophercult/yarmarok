package storage

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/kaznasho/yarmarok/service"
)

// RaffleStorage returns a storage for raffles.
func (s *FirestoreOrganizerStorage) RaffleStorage(organizerID string) service.RaffleStorage {
	return NewFirestoreRaffleStorage(s.firestoreClient.Doc(organizerID).Collection(raffleCollection), organizerID)
}

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
}

// Create creates a new raffle.
func (s *FirestoreRaffleStorage) Create(raf *service.Raffle) error {
	doc, err := s.Get(raf.ID)
	if !isNotFound(err) {
		if err != nil {
			return fmt.Errorf("check existence: %w", err)
		}

		if doc != nil {
			return service.ErrRaffleAlreadyExists
		}
	}

	raf.OrganizerID = s.organizerID

	_, err = s.firestoreClient.Doc(raf.ID).Set(context.Background(), raf)

	return err
}

// Get returns a raffle with the given id.
func (s *FirestoreRaffleStorage) Get(id string) (*service.Raffle, error) {
	doc, err := s.firestoreClient.Doc(id).Get(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get document: %w", err)
	}

	var raf service.Raffle
	err = doc.DataTo(&raf)
	if err != nil {
		return nil, fmt.Errorf("decode document: %w", err)
	}

	return &raf, nil
}

// GetAll returns all raffles per organizer.
func (s *FirestoreRaffleStorage) GetAll() ([]service.Raffle, error) {
	docs, err := s.firestoreClient.Documents(context.Background()).GetAll()
	if err != nil {
		return nil, fmt.Errorf("get documents: %w", err)
	}

	var raffles []service.Raffle
	for _, doc := range docs {
		var raf service.Raffle
		err = doc.DataTo(&raf)
		if err != nil {
			return nil, fmt.Errorf("decode document: %w", err)
		}

		raffles = append(raffles, raf)
	}

	return raffles, nil
}
