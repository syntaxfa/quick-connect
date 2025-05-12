package outboxpsq

import (
	"context"
	"time"

	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

const queryRemoveRecordsBeforeDatetime = `DELETE FROM outbox
WHERE created_on < $1;`

func (d *DB) RemoveRecordsBeforeDatetime(expiryTime time.Time) error {
	const op = "outbox.repository.remove.RemoveRecordsBeforeDatetime"

	if _, eErr := d.conn.Exec(context.Background(), queryRemoveRecordsBeforeDatetime, expiryTime); eErr != nil {
		return richerror.New(op).WithWrapError(eErr).WithKind(richerror.KindUnexpected)
	}

	return nil
}
