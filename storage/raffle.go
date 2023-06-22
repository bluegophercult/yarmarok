package storage

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/kaznasho/yarmarok/service"
)

// RaffleStorage returns a storage for raffles.
func (os *FirestoreOrganizerStorage) RaffleStorage(organizerID string) service.RaffleStorage {
	return NewFirestoreRaffleStorage(os.firestoreClient.Doc(organizerID).Collection(raffleCollection), organizerID)
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
func (ys *FirestoreRaffleStorage) Create(y *service.Raffle) error {
	doc, err := ys.Get(y.ID)
	if !isNotFound(err) {
		if err != nil {
			return fmt.Errorf("check existence: %w", err)
		}

		if doc != nil {
			return service.ErrRaffleAlreadyExists
		}
	}

	y.OrganizerID = ys.organizerID

	_, err = ys.firestoreClient.Doc(y.ID).Set(context.Background(), y)

	return err
}

// Get returns a raffle with the given id.
func (ys *FirestoreRaffleStorage) Get(id string) (*service.Raffle, error) {
	doc, err := ys.firestoreClient.Doc(id).Get(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get document: %w", err)
	}

	var y service.Raffle
	err = doc.DataTo(&y)
	if err != nil {
		return nil, fmt.Errorf("decode document: %w", err)
	}

	return &y, nil
}

// GetAll returns all raffles per organizer.
func (ys *FirestoreRaffleStorage) GetAll() ([]service.Raffle, error) {
	docs, err := ys.firestoreClient.Documents(context.Background()).GetAll()
	if err != nil {
		return nil, fmt.Errorf("get documents: %w", err)
	}

	var raffles []service.Raffle
	for _, doc := range docs {
		var y service.Raffle
		err = doc.DataTo(&y)
		if err != nil {
			return nil, fmt.Errorf("decode document: %w", err)
		}

		raffles = append(raffles, y)
	}

	return raffles, nil
}
