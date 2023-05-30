package yarmarok

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	userCollection = "users"
)

var (
	// ErrUserAlreadyExists is returned when a user already exists.
	ErrUserAlreadyExists = errors.New("user already exists")
)

// User represents a user of the application.
type User struct {
	ID string
}

// UserStorage is a storage for users.
type UserStorage interface {
	Create(User) error
	Exists(id string) (bool, error)
}

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
func (us *FirestoreUserStorage) Create(u User) error {
	exists, err := us.Exists(u.ID)
	if err != nil {
		return fmt.Errorf("check user exists: %w", err)
	}

	if exists {
		return ErrUserAlreadyExists
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
