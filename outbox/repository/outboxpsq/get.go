package outboxpsq

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
	"errors"

	"github.com/syntaxfa/quick-connect/outbox"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

const queryGetRecordsByLockID = `SELECT
id, data, state, created_on, locked_by, locked_on, processed_on, number_of_attempts, last_attempted_on, error
FROM outbox
WHERE locked_by=$1;
`

func (d *DB) GetRecordsByLockID(lockID string) ([]outbox.Record, error) {
	const op = "outbox.repository.outboxpsq.get.GetRecordsByLockID"

	rows, qErr := d.conn.Query(context.Background(), queryGetRecordsByLockID, lockID)
	if qErr != nil {
		return nil, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}
	defer rows.Close()

	var messages []outbox.Record

	for rows.Next() {
		var rec outbox.Record
		var data []byte

		sErr := rows.Scan(&rec.ID, &data, &rec.State, &rec.CreatedOn, &rec.LockID, &rec.LockedOn,
			&rec.ProcessedOn, &rec.NumberOfAttempts, &rec.LastAttemptOn, &rec.Error)
		if sErr != nil {
			if errors.Is(sErr, sql.ErrNoRows) {
				return messages, nil
			}

			return nil, richerror.New(op).WithWrapError(sErr).WithKind(richerror.KindUnexpected)
		}
		if dErr := gob.NewDecoder(bytes.NewReader(data)).Decode(&rec.Message); dErr != nil {
			return nil, richerror.New(op).WithWrapError(dErr).WithKind(richerror.KindUnexpected)
		}

		messages = append(messages, rec)
	}

	if err := rows.Err(); err != nil {
		return messages, richerror.New(op).WithWrapError(err).WithKind(richerror.KindUnexpected)
	}

	return messages, nil
}
