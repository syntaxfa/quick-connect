package outbox

import (
	"database/sql"
	"github.com/syntaxfa/quick-connect/outbox/internal/time"
	"github.com/syntaxfa/quick-connect/outbox/internal/uuid"
)

type Publisher struct {
	store Store
	time  time.Provider
	uuid  uuid.Provider
}

func NewPublisher(store Store) Publisher {
	return Publisher{
		store: store,
		time:  time.NewTimeProvider(),
		uuid:  uuid.NewUUIDProvider(),
	}
}

func (p Publisher) Send(message Message, tx *sql.Tx) error {
	newID := p.uuid.NewUUID()

	record := Record{
		ID:               newID,
		Message:          message,
		State:            StatePendingDelivery,
		CreatedOn:        p.time.Now(),
		LockID:           nil,
		LockedOn:         nil,
		ProcessedOn:      nil,
		NumberOfAttempts: 0,
		LastAttemptOn:    nil,
		Error:            nil,
	}

	return p.store.AddRecordTx(record, tx)
}
