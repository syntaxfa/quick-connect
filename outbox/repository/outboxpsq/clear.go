package outboxpsq

import (
	"context"
	"time"

	"github.com/syntaxfa/quick-connect/adapter/postgres"
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

	stmt, pErr := d.conn.PrepareStatement(context.Background(), postgres.StatementClearLocksWithDurationBeforeDate, queryClearLocksWithDurationBeforeDate) //nolint:sqlclosecheck // finally closed, but not here
	if pErr != nil {
		return richerror.New(op).WithWrapError(pErr).WithKind(richerror.KindUnexpected)
	}

	_, eErr := stmt.Exec(time)
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

	stmt, pErr := d.conn.PrepareStatement(context.Background(), postgres.StatementClearLocksByLockID, queryClearLocksByLockID) //nolint:sqlclosecheck // finally closed, but not here
	if pErr != nil {
		return richerror.New(op).WithWrapError(pErr).WithKind(richerror.KindUnexpected)
	}

	if _, eErr := stmt.Exec(lockID); eErr != nil {
		return richerror.New(op).WithWrapError(eErr).WithKind(richerror.KindUnexpected)
	}

	return nil
}
