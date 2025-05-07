package outbox

import (
	time2 "time"

	"github.com/syntaxfa/quick-connect/outbox/internal/time"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

type recordUnlocker struct {
	store                  Store
	time                   time.Provider
	MaxLockTimeDurationMin time2.Duration
}

func newRecordUnlocker(store Store, maxLockTimeDurationMin time2.Duration) recordUnlocker {
	return recordUnlocker{
		store:                  store,
		time:                   time.NewTimeProvider(),
		MaxLockTimeDurationMin: maxLockTimeDurationMin,
	}
}

func (d recordUnlocker) UnlockExpiredMessages() error {
	const op = "outbox.record_unlocker.UnlockExpiredMessages"

	expiryTime := d.time.Now().UTC().Add(-d.MaxLockTimeDurationMin)

	if cErr := d.store.ClearLocksWithDurationBeforeDate(expiryTime); cErr != nil {
		return richerror.New(op).WithWrapError(cErr).WithKind(richerror.KindUnexpected)
	}

	return nil
}
