package storage

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"

	"github.com/kaznasho/yarmarok/service"
)

// FirestoreParticipantStorage is a storage for raffles based on Firestore.
type FirestoreParticipantStorage struct {
	raffleID        string
	firestoreClient *firestore.CollectionRef
}

// NewFirestoreParticipantStorage creates a new FirestoreParticipantStorage.
func NewFirestoreParticipantStorage(client *firestore.CollectionRef, raffleID string) *FirestoreParticipantStorage {
	return &FirestoreParticipantStorage{
		raffleID:        raffleID,
		firestoreClient: client,
	}
}

// Create creates a new participant.
func (ps *FirestoreParticipantStorage) Create(p *service.Participant) error {
	exists, err := ps.Exists(p.ID)
	if err != nil {
		return fmt.Errorf("check participant exists: %w", err)
	}

	if exists {
		return service.ErrParticipantAlreadyExists
	}

	if _, err := ps.firestoreClient.Doc(p.ID).Set(context.Background(), p); err != nil {
		return fmt.Errorf("create participant: %w", err)
	}

	return nil
}

// Get returns a participant with the given ID.
func (ps *FirestoreParticipantStorage) Get(id string) (*service.Participant, error) {
	doc, err := ps.firestoreClient.Doc(id).Get(context.Background())
	if err != nil {
		if isNotFound(err) {
			return nil, service.ErrParticipantNotFound
		}
		return nil, fmt.Errorf("get participant: %w", err)
	}

	var p service.Participant
	if err = doc.DataTo(&p); err != nil {
		return nil, fmt.Errorf("decode participant: %w", err)
	}

	return &p, nil
}

// Update updates an existing participant.
func (ps *FirestoreParticipantStorage) Update(p *service.Participant) error {
	exists, err := ps.Exists(p.ID)
	if err != nil {
		return fmt.Errorf("check participant exists: %w", err)
	}

	if !exists {
		return service.ErrParticipantNotFound
	}

	if _, err := ps.firestoreClient.Doc(p.ID).Set(context.Background(), p); err != nil {
		return fmt.Errorf("update participant: %w", err)
	}

	return nil
}

// Delete deletes a participant.
func (ps *FirestoreParticipantStorage) Delete(id string) error {
	exists, err := ps.Exists(id)
	if err != nil {
		return fmt.Errorf("check participant exists: %w", err)
	}

	if !exists {
		return service.ErrParticipantNotFound
	}

	if _, err := ps.firestoreClient.Doc(id).Delete(context.Background()); err != nil {
		return fmt.Errorf("delete participant: %w", err)
	}

	return nil
}

// GetAll returns a list of all participants.
func (ps *FirestoreParticipantStorage) GetAll() ([]service.Participant, error) {
	docs, err := ps.firestoreClient.Documents(context.Background()).GetAll()
	if err != nil {
		return nil, fmt.Errorf("get participant: %w", err)
	}

	participants := make([]service.Participant, 0, len(docs))
	for _, doc := range docs {
		var p service.Participant
		if err = doc.DataTo(&p); err != nil {
			return nil, fmt.Errorf("decode participants: %w", err)
		}

		participants = append(participants, p)
	}

	return participants, nil
}

// Exists check if Participant exists
func (ps *FirestoreParticipantStorage) Exists(id string) (bool, error) {
	doc, err := ps.firestoreClient.Doc(id).Get(context.Background())
	if isNotFound(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return doc.Exists(), nil
}
