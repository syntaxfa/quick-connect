package outboxpsq

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"

	"github.com/syntaxfa/quick-connect/outbox"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

const queryAddRecordTx = `INSERT INTO outbox
(id, data, state, created_on, locked_by, locked_on, processed_on, number_of_attempts, last_attempted_on, error)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`

func (d *DB) AddRecordTX(rec outbox.Record, _ *sql.Tx) error {
	const op = "outbox.repository.outboxpsq.add.AddRecordTX"

	msgBuf := new(bytes.Buffer)
	msgEnc := gob.NewEncoder(msgBuf)
	if eErr := msgEnc.Encode(rec.Message); eErr != nil {
		return richerror.New(op).WithWrapError(eErr).WithKind(richerror.KindUnexpected)
	}

	if _, eErr := d.conn.Exec(context.Background(), queryAddRecordTx, rec.ID, msgBuf.Bytes(), rec.State, rec.CreatedOn, rec.LockID, rec.LockedOn, rec.ProcessedOn,
		rec.NumberOfAttempts, rec.LastAttemptOn, rec.Error); eErr != nil {
		return richerror.New(op).WithWrapError(eErr).WithKind(richerror.KindUnexpected)
	}

	return nil
}
