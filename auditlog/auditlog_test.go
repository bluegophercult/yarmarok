package auditlog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAuditLogRecord(t *testing.T) {
	actor := "actor"
	action := "action"
	context := struct {
		foo string
		bar int
	}{
		foo: "foo",
		bar: 1,
	}

	record := NewAuditLogRecord(actor, action, context)

	assert.Equal(t, actor, record.ActorID)
	assert.Equal(t, action, record.Action)
	assert.Equal(t, context, record.Context)

	assert.NotEmpty(t, record.ID)
	assert.NotEmpty(t, record.CreatedAt)
}
