package storage

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/kaznasho/yarmarok/service"
)

// YarmarokStorage returns a storage for yarmaroks.
func (us *FirestoreUserStorage) YarmarokStorage(userID string) service.YarmarokStorage {
	return NewFirestoreYarmarokStorage(us.firestoreClient.Doc(userID).Collection(yarmarokCollection), userID)
}

// NewFirestoreYarmarokStorage creates a new FirestoreYarmarokStorage.
func NewFirestoreYarmarokStorage(client *firestore.CollectionRef, userID string) *FirestoreYarmarokStorage {
	return &FirestoreYarmarokStorage{
		userID:          userID,
		firestoreClient: client,
	}
}

// FirestoreYarmarokStorage is a storage for yarmaroks based on Firestore.
type FirestoreYarmarokStorage struct {
	userID          string
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

	y.UserID = ys.userID

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
