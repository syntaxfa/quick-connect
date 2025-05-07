package outbox

import (
	"time"

	"github.com/google/uuid"
)

type Record struct {
	ID               uuid.UUID
	Message          Message
	State            RecordState
	CreatedOn        time.Time
	LockID           *string
	LockedOn         *time.Time
	ProcessedOn      *time.Time
	NumberOfAttempts int
	LastAttemptOn    *time.Time
	Error            *string
}

type Message struct {
	Key     string
	Headers map[string]string
	Body    []byte
	Topic   string
}

type RecordState int

const (
	StatePendingDelivery RecordState = iota + 1
	StateDelivered
	StateMaxAttemptsReached
)
