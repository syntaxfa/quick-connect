package outboxpsq

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"github.com/syntaxfa/quick-connect/outbox"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

const queryUpdateRecordLockeByState = `UPDATE outbox
SET
	locked_by=$1,
	locked_on=$2,
WHERE state=$3;
`

func (d *DB) UpdateRecordLockByState(lockID string, lockedOn time.Time, state outbox.RecordState) error {
	const op = "outbox.repository.outboxpsq.update.UpdateRecordLockByState"

	if _, eErr := d.conn.Exec(context.Background(), queryUpdateRecordLockeByState, lockID, lockedOn, state); eErr != nil {
		return richerror.New(op).WithWrapError(eErr).WithKind(richerror.KindUnexpected)
	}

	return nil
}

const queryUpdateRecordByID = `UPDATE outbox
SET
    data=$1,
    state=$2,
    created_on=$3,
    locked_by=$4,
    locked_on=$5,
    processed_on=$6,
    number_of_attempts=$6,
    last_attempted_on=$7,
    error=$8,
WHERE id = $9
`

func (d *DB) UpdateRecordByID(rec outbox.Record) error {
	const op = "outbox.repository.outboxpsq.update.UpdateRecordByID"

	msgData := new(bytes.Buffer)
	enc := gob.NewEncoder(msgData)
	if eErr := enc.Encode(rec.Message); eErr != nil {
		return richerror.New(op).WithWrapError(eErr).WithKind(richerror.KindUnexpected)
	}

	_, eErr := d.conn.Exec(context.Background(), queryUpdateRecordByID, msgData, rec.State, rec.CreatedOn, rec.LockID, rec.LockedOn,
		rec.ProcessedOn, rec.NumberOfAttempts, rec.LastAttemptOn, rec.Error, rec.ID)
	if eErr != nil {
		return richerror.New(op).WithWrapError(eErr).WithKind(richerror.KindUnexpected)
	}

	return nil
}
