package storage

import (
	"cloud.google.com/go/firestore"

	"github.com/kaznasho/yarmarok/auditlog"
)

// FirestoreAuditLogStorage is a storage for audit logs based on Firestore.
type FirestoreAuditLogStorage struct {
	*StorageBase[auditlog.AuditLogRecord]
}

// NewFirestoreAuditLogStorage creates a new FirestoreAuditLogStorage.
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
