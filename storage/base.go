package storage

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kaznasho/yarmarok/service"

	"cloud.google.com/go/firestore"
)

// Storable is a type parameter constraint for all storable items.
type Storable interface {
	service.Raffle | service.Prize | service.Participant | service.Organizer | service.Donation
}

// IDExtractor is a typed function that extracts an ID from the item it serves.
type IDExtractor[Item Storable] func(*Item) string

// StorageBase is a base with common functionality for all storages.
type StorageBase[Item Storable] struct {
	collectionReference *firestore.CollectionRef
	extractID           IDExtractor[Item]
}

// NewStorageBase creates a new StorageBase.
func NewStorageBase[Item Storable](collectionReference *firestore.CollectionRef, idExtractor IDExtractor[Item]) *StorageBase[Item] {
	return &StorageBase[Item]{
		collectionReference: collectionReference,
		extractID:           idExtractor,
	}
}

// Create creates a new item.
func (sb *StorageBase[Item]) Create(item *Item) error {
	id := sb.extractID(item)
	exists, err := sb.Exists(id)
	if err != nil {
		return fmt.Errorf("check item exists: %w", err)
	}

	if exists {
		return service.ErrAlreadyExists
	}

	_, err = sb.collectionReference.Doc(id).Set(context.Background(), item)
	if err != nil {
		return fmt.Errorf("create item: %w", err)
	}

	return nil
}

// Get returns an item with the given ID.
func (sb *StorageBase[Item]) Get(id string) (*Item, error) {
	doc, err := sb.collectionReference.Doc(id).Get(context.Background())
	if err != nil {
		if isNotFound(err) {
			return nil, service.ErrNotFound
		}
		return nil, fmt.Errorf("get item: %w", err)
	}

	var i Item
	if err := doc.DataTo(&i); err != nil {
		return nil, fmt.Errorf("decode item: %w", err)
	}

	return &i, nil
}

// Update replaces an item with the given ID with the given item.
func (sb *StorageBase[Item]) Update(item *Item) error {
	id := sb.extractID(item)
	exists, err := sb.Exists(id)
	if err != nil {
		return fmt.Errorf("check item exists: %w", err)
	}

	if !exists {
		return service.ErrNotFound
	}

	_, err = sb.collectionReference.Doc(id).Set(context.Background(), item)
	if err != nil {
		return fmt.Errorf("create item: %w", err)
	}

	return nil
}

// GetAll returns all items in the collection.
func (sb *StorageBase[Item]) GetAll() ([]Item, error) {
	docs, err := sb.collectionReference.Documents(context.Background()).GetAll()
	if err != nil {
		return nil, fmt.Errorf("get all items: %w", err)
	}

	// TODO: ErrNotFound if no docs.

	items := make([]Item, 0, len(docs))

	for _, doc := range docs {
		var item Item
		if err = doc.DataTo(&item); err != nil {
			return nil, fmt.Errorf("decode items: %w", err)
		}

		items = append(items, item)
	}

	return items, nil
}

// Delete deletes an item with the given ID.
func (sb *StorageBase[Item]) Delete(id string) error {
	exists, err := sb.Exists(id)
	if err != nil {
		return fmt.Errorf("check item exists: %w", err)
	}

	if !exists {
		return service.ErrNotFound
	}

	_, err = sb.collectionReference.Doc(id).Delete(context.Background())
	if err != nil {
		return fmt.Errorf("delete item: %w", err)
	}

	return nil
}

// Exists checks if an item with the given ID exists.
func (sb *StorageBase[Item]) Exists(id string) (bool, error) {
	doc, err := sb.collectionReference.Doc(id).Get(context.Background())
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

func (sb *StorageBase[Item]) Query(query service.Query) ([]Item, error) {
	q := sb.collectionReference.Query

	for f := query.Filter; f != nil; f = f.Next {
		q = q.Where(f.Field, toOp(f.Operator), f.Value)
	}

	if query.OrderBy != nil {
		q = q.OrderBy(query.OrderBy.Field, toDir(query.OrderBy.Direction))
	}

	if query.Limit > 0 {
		q = q.Limit(query.Limit)
	}

	if query.Offset > 0 {
		q = q.Offset(query.Offset)
	}

	docs, err := q.Documents(context.Background()).GetAll()
	if err != nil {
		return nil, fmt.Errorf("querying items: %w", err)
	}

	l := len(docs)
	if l == 0 {
		return nil, service.ErrNotFound
	}

	items := make([]Item, l)

	for i := range docs {
		var item Item
		if err := docs[i].DataTo(&item); err != nil {
			return nil, fmt.Errorf("decode items: %w", err)
		}

		items[i] = item
	}

	return items, nil

}

func toOp(op service.Operator) string {
	switch op {
	case service.LT:
		return "<"
	case service.LTE:
		return "<="
	case service.EQ:
		return "=="
	case service.GT:
		return ">"
	case service.GTE:
		return ">="
	case service.NE:
		return "!="
	case service.IN:
		return "in"
	case service.NI:
		return "not_in"
	case service.CN:
		return "array_contains"
	case service.CA:
		return "array_contains_any"
	default:
		return ""
	}
}

func toDir(op service.Direction) firestore.Direction {
	switch op {
	case service.ASC:
		return firestore.Asc
	case service.DESC:
		return firestore.Desc
	default:
		return 0
	}
}
