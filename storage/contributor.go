package storage

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/kaznasho/yarmarok/service"
)

// FirestoreContributorStorage is a storage for raffles based on Firestore.
type FirestoreContributorStorage struct {
	raffleID        string
	firestoreClient *firestore.CollectionRef
}

// ContributorStorage returns a contributor storage.
func (s *FirestoreRaffleStorage) ContributorStorage(raffleID string) service.ContributorStorage {
	return NewFirestoreContributorStorage(s.firestoreClient.Doc(raffleID).Collection(contributorCollection), raffleID)
}

// NewFirestoreContributorStorage creates a new FirestoreContributorStorage.
func NewFirestoreContributorStorage(client *firestore.CollectionRef, raffleID string) *FirestoreContributorStorage {
	return &FirestoreContributorStorage{
		raffleID:        raffleID,
		firestoreClient: client,
	}
}

// Create creates a new contributor.
func (s *FirestoreContributorStorage) Create(ctb *service.Contributor) error {
	exists, err := s.Exists(ctb.ID)
	if err != nil {
		return fmt.Errorf("check contributor exists: %w", err)
	}

	if exists {
		return service.ErrContributorAlreadyExists
	}

	if _, err := s.firestoreClient.Doc(ctb.ID).Set(context.Background(), ctb); err != nil {
		return fmt.Errorf("create contributor: %w", err)
	}

	return nil
}

// Get returns a contributor with the given ID.
func (s *FirestoreContributorStorage) Get(id string) (*service.Contributor, error) {
	doc, err := s.firestoreClient.Doc(id).Get(context.Background())
	if err != nil {
		if isNotFound(err) {
			return nil, service.ErrContributorNotFound
		}
		return nil, fmt.Errorf("get contributor: %w", err)
	}

	var ctb service.Contributor
	if err = doc.DataTo(&ctb); err != nil {
		return nil, fmt.Errorf("decode contributor: %w", err)
	}

	return &ctb, nil
}

// Update updates an existing contributor.
func (s *FirestoreContributorStorage) Update(ctb *service.Contributor) error {
	exists, err := s.Exists(ctb.ID)
	if err != nil {
		return fmt.Errorf("check contributor exists: %w", err)
	}

	if !exists {
		return service.ErrContributorNotFound
	}

	if _, err := s.firestoreClient.Doc(ctb.ID).Set(context.Background(), ctb); err != nil {
		return fmt.Errorf("update contributor: %w", err)
	}

	return nil
}

// Delete deletes a contributor.
func (s *FirestoreContributorStorage) Delete(id string) error {
	exists, err := s.Exists(id)
	if err != nil {
		return fmt.Errorf("check contributor exists: %w", err)
	}

	if !exists {
		return service.ErrContributorNotFound
	}

	if _, err := s.firestoreClient.Doc(id).Delete(context.Background()); err != nil {
		return fmt.Errorf("delete contributor: %w", err)
	}

	return nil
}

// GetAll returns a list of all contributors.
func (s *FirestoreContributorStorage) GetAll() ([]service.Contributor, error) {
	docs, err := s.firestoreClient.Documents(context.Background()).GetAll()
	if err != nil {
		return nil, fmt.Errorf("get contributor: %w", err)
	}

	contributors := make([]service.Contributor, 0, len(docs))
	for _, doc := range docs {
		var ctb service.Contributor
		if err = doc.DataTo(&ctb); err != nil {
			return nil, fmt.Errorf("decode contributors: %w", err)
		}

		contributors = append(contributors, ctb)
	}

	return contributors, nil
}

func (s *FirestoreContributorStorage) Exists(id string) (bool, error) {
	doc, err := s.firestoreClient.Doc(id).Get(context.Background())
	if isNotFound(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return doc.Exists(), nil
}
