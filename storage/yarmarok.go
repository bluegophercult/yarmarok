package storage

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/kaznasho/yarmarok/service"
)

// YarmarokStorage returns a storage for yarmaroks.
func (os *FirestoreOrganizerStorage) YarmarokStorage(organizerID string) service.YarmarokStorage {
	return NewFirestoreYarmarokStorage(os.firestoreClient.Doc(organizerID).Collection(yarmarokCollection), organizerID)
}

// NewFirestoreYarmarokStorage creates a new FirestoreYarmarokStorage.
func NewFirestoreYarmarokStorage(client *firestore.CollectionRef, organizerID string) *FirestoreYarmarokStorage {
	return &FirestoreYarmarokStorage{
		organizerID:     organizerID,
		firestoreClient: client,
	}
}

// FirestoreYarmarokStorage is a storage for yarmaroks based on Firestore.
type FirestoreYarmarokStorage struct {
	organizerID     string
	firestoreClient *firestore.CollectionRef
}

// Create creates a new yarmarok.
func (ys *FirestoreYarmarokStorage) Create(y *service.Yarmarok) error {
	doc, err := ys.Get(y.ID)
	if !isNotFound(err) {
		if err != nil {
			return fmt.Errorf("check existence: %w", err)
		}

		if doc != nil {
			return service.ErrYarmarokAlreadyExists
		}
	}

	y.OrganizerID = ys.organizerID

	_, err = ys.firestoreClient.Doc(y.ID).Set(context.Background(), y)

	return err
}

// Get returns a yarmarok with the given id.
func (ys *FirestoreYarmarokStorage) Get(id string) (*service.Yarmarok, error) {
	doc, err := ys.firestoreClient.Doc(id).Get(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get document: %w", err)
	}

	var y service.Yarmarok
	err = doc.DataTo(&y)
	if err != nil {
		return nil, fmt.Errorf("decode document: %w", err)
	}

	return &y, nil
}

// GetAll returns all yarmaroks per organizer.
func (ys *FirestoreYarmarokStorage) GetAll() ([]service.Yarmarok, error) {
	docs, err := ys.firestoreClient.Documents(context.Background()).GetAll()
	if err != nil {
		return nil, fmt.Errorf("get documents: %w", err)
	}

	var yarmaroks []service.Yarmarok
	for _, doc := range docs {
		var y service.Yarmarok
		err = doc.DataTo(&y)
		if err != nil {
			return nil, fmt.Errorf("decode document: %w", err)
		}

		yarmaroks = append(yarmaroks, y)
	}

	return yarmaroks, nil
}
