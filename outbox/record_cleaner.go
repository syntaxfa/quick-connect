package outbox

import (
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	time2 "time"

	"github.com/syntaxfa/quick-connect/outbox/internal/time"
)

type recordCleaner struct {
	store             Store
	time              time.Provider
	MaxRecordLifetime time2.Duration
}

func newRecordCleaner(store Store, maxRecordLifetime time2.Duration) recordCleaner {
	return recordCleaner{
		store:             store,
		time:              time.NewTimeProvider(),
		MaxRecordLifetime: maxRecordLifetime,
	}
}

func (d recordCleaner) RemoveExpiredMessages() error {
	const op = "outbox.record_cleaner.RemoveExpiredMessages"

	expiryTime := d.time.Now().UTC().Add(-d.MaxRecordLifetime)

	if rErr := d.store.RemoveRecordsBeforeDatetime(expiryTime); rErr != nil {
		return richerror.New(op).WithWrapError(rErr).WithKind(richerror.KindUnexpected)
	}

	return nil
}
