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
	userCollection     = "users"
	yarmarokCollection = "yarmaroks"
)

// FirestoreUserStorage is a storage for users based on Firestore.
type FirestoreUserStorage struct {
	firestoreClient *firestore.CollectionRef
}

// NewFirestoreUserStorage creates a new FirestoreUserStorage.
func NewFirestoreUserStorage(client *firestore.Client) *FirestoreUserStorage {
	return &FirestoreUserStorage{
		firestoreClient: client.Collection(userCollection),
	}
}

// Create creates a new user.
func (us *FirestoreUserStorage) Create(u service.User) error {
	exists, err := us.Exists(u.ID)
	if err != nil {
		return fmt.Errorf("check user exists: %w", err)
	}

	if exists {
		return service.ErrUserAlreadyExists
	}

	_, err = us.firestoreClient.Doc(u.ID).Set(context.Background(), u)
	return err
}

// Exists checks if a user with the given id exists.
func (us *FirestoreUserStorage) Exists(id string) (bool, error) {
	doc, err := us.firestoreClient.Doc(id).Get(context.Background())
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
