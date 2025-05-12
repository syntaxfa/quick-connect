package outboxpsq

import (
	"context"
	"time"

	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

const queryClearLocksWithDurationBeforeDate = `UPDATE outbox
SET
	locked_by=NULL,
	locked_on=NULL,
WHERE locked_on < $1;
`

func (d *DB) ClearLocksWithDurationBeforeDate(time time.Time) error {
	const op = "outbox.repository.outboxpsq.remove.ClearLocksWithDurationBeforeDate"

	_, eErr := d.conn.Exec(context.Background(), queryClearLocksWithDurationBeforeDate, time)
	if eErr != nil {
		return richerror.New(op).WithWrapError(eErr).WithKind(richerror.KindUnexpected)
	}

	return nil
}

const queryClearLocksByLockID = `UPDATE outbox
SET
    locked_by=NULL,
    locked_on=NULL,
WHERE locked_by = $1;
`

func (d *DB) ClearLocksByLockID(lockID string) error {
	const op = "outbox.repository.outboxpsq.remove.ClearLocksByLockID"

	if _, eErr := d.conn.Exec(context.Background(), queryClearLocksByLockID, lockID); eErr != nil {
		return richerror.New(op).WithWrapError(eErr).WithKind(richerror.KindUnexpected)
	}

	return nil
}
