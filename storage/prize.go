package storage

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"

	"github.com/kaznasho/yarmarok/service"
)

// FirestorePrizeStorage is a storage for prizes based on Firestore.
type FirestorePrizeStorage struct {
	yarmarokID      string
	firestoreClient *firestore.CollectionRef
}

// PrizeStorage returns a prize storage.
func (ys *FirestoreYarmarokStorage) PrizeStorage(yarmarokID string) service.PrizeStorage {
	return NewFirestorePrizeStorage(ys.firestoreClient.Doc(yarmarokID).Collection(prizeCollection), yarmarokID)
}

// NewFirestorePrizeStorage creates a new FirestorePrizeStorage.
func NewFirestorePrizeStorage(client *firestore.CollectionRef, yarmarokID string) *FirestorePrizeStorage {
	return &FirestorePrizeStorage{
		yarmarokID:      yarmarokID,
		firestoreClient: client,
	}
}

// Create creates a new prize.
func (ps *FirestorePrizeStorage) Create(p *service.Prize) error {
	exists, err := ps.Exists(p.ID)
	if err != nil {
		return fmt.Errorf("check prize exists: %w", err)
	}

	if exists {
		return service.ErrPrizeAlreadyExists
	}

	if _, err := ps.firestoreClient.Doc(p.ID).Set(context.Background(), p); err != nil {
		return fmt.Errorf("create prize: %w", err)
	}

	return nil
}

// Get returns a prize with the given ID.
func (ps *FirestorePrizeStorage) Get(id string) (*service.Prize, error) {
	doc, err := ps.firestoreClient.Doc(id).Get(context.Background())
	if err != nil {
		if isNotFound(err) {
			return nil, service.ErrPrizeNotFound
		}
		return nil, fmt.Errorf("get prize: %w", err)
	}

	var p service.Prize
	if err := doc.DataTo(&p); err != nil {
		return nil, fmt.Errorf("decode prize: %w", err)
	}

	return &p, nil
}

// Update updates an existing prize.
func (ps *FirestorePrizeStorage) Update(p *service.Prize) error {
	exists, err := ps.Exists(p.ID)
	if err != nil {
		return fmt.Errorf("check prize exists: %w", err)
	}

	if !exists {
		return service.ErrPrizeNotFound
	}

	if _, err := ps.firestoreClient.Doc(p.ID).Set(context.Background(), p); err != nil {
		return fmt.Errorf("update prize: %w", err)
	}

	return nil
}

// GetAll returns a list of all prizes.
func (ps *FirestorePrizeStorage) GetAll() ([]service.Prize, error) {
	docs, err := ps.firestoreClient.Documents(context.Background()).GetAll()
	if err != nil {
		return nil, fmt.Errorf("get all prizes: %w", err)
	}

	prizes := make([]service.Prize, 0, len(docs))
	for _, doc := range docs {
		var p service.Prize
		if err = doc.DataTo(&p); err != nil {
			return nil, fmt.Errorf("decode prizes: %w", err)
		}

		prizes = append(prizes, p)
	}

	return prizes, nil
}

// Delete deletes a prize.
func (ps *FirestorePrizeStorage) Delete(id string) error {
	exists, err := ps.Exists(id)
	if err != nil {
		return fmt.Errorf("check prize exists: %w", err)
	}

	if !exists {
		return service.ErrPrizeNotFound
	}

	if _, err := ps.firestoreClient.Doc(id).Delete(context.Background()); err != nil {
		return fmt.Errorf("delete prize: %w", err)
	}

	return nil
}

// Exists check if prize exists
func (ps *FirestorePrizeStorage) Exists(id string) (bool, error) {
	doc, err := ps.firestoreClient.Doc(id).Get(context.Background())
	if isNotFound(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return doc.Exists(), nil
}
