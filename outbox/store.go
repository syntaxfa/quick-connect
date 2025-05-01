package outbox

import (
	"database/sql"
	"time"
)

type Store interface {
	AddRecordTx(record Record, tx *sql.Tx) error
	GetRecordsByLockID(lockID string) ([]Record, error)
	UpdateRecordLockByState(lockID string, lockedOn time.Time, state RecordState) error
	UpdateRecordByID(message Record) error
	ClearLocksWithDurationBeforeDate(time time.Time) error
	ClearLocksByLockID(lockID string) error
	RemoveRecordsBeforeDatetime(expiryTime time.Time) error
}
