package outboxpsq

import (
	"context"
	"time"

	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

const queryRemoveRecordsBeforeDatetime = `DELETE FROM outbox
WHERE created_on < $1;`

func (d *DB) RemoveRecordsBeforeDatetime(expiryTime time.Time) error {
	const op = "outbox.repository.remove.RemoveRecordsBeforeDatetime"

	stmt, pErr := d.conn.PrepareStatement(context.Background(), postgres.StatementRemoveRecordsBeforeDatetime, queryRemoveRecordsBeforeDatetime) //nolint:sqlclosecheck // finally closed, but not here
	if pErr != nil {
		return richerror.New(op).WithWrapError(pErr).WithKind(richerror.KindUnexpected)
	}

	if _, eErr := stmt.Exec(expiryTime); eErr != nil {
		return richerror.New(op).WithWrapError(eErr).WithKind(richerror.KindUnexpected)
	}

	return nil
}
