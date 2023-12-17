package storage

import (
	"cloud.google.com/go/firestore"

	"github.com/kaznasho/yarmarok/auditlog"
)

// FirestoreDonationStorage is a storage for donation based on Firestore.
type FirestoreAuditLogStorage struct {
	*StorageBase[auditlog.AuditLogRecord]
}

// NewFirestoreDonationStorage creates a new FirestoreDonationStorage.
func NewFirestoreAuditLogStorage(
	client *firestore.CollectionRef,
) *FirestoreAuditLogStorage {
	recordIDExtractor := IDExtractor[auditlog.AuditLogRecord](
		func(r *auditlog.AuditLogRecord) string {
			return r.ID
		},
	)

	return &FirestoreAuditLogStorage{
		StorageBase: NewStorageBase[auditlog.AuditLogRecord](
			client,
			recordIDExtractor,
		),
	}
}
