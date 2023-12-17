package auditlog

import (
	"time"

	"github.com/google/uuid"
)

var timeNow = time.Now().UTC

var stringUUID = uuid.New().String

//go:generate mockgen -destination=mock_auditlog_storage_test.go -package=auditlog github.com/kaznasho/yarmarok/auditlog AuditLogStorage
type AuditLogStorage interface {
	Create(*AuditLogRecord) error
	GetAll() ([]AuditLogRecord, error)
}

type AuditLogRecord struct {
	ID        string    `json:"id"`
	ActorID   string    `json:"actorId"`
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"createdAt"`
	Context   any       `json:"context"`
}

func NewAuditLogRecord(actor, action string, context any) *AuditLogRecord {
	return &AuditLogRecord{
		ID:        stringUUID(),
		ActorID:   actor,
		Action:    action,
		CreatedAt: timeNow(),
		Context:   context,
	}
}
