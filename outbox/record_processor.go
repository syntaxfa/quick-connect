package outbox

import (
	"fmt"
	"log/slog"

	"github.com/syntaxfa/quick-connect/outbox/internal/time"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

type defaultRecordProcessor struct {
	messageBroker MessageBroker
	store         Store
	time          time.Provider
	machineID     string
	retrialPolicy RetrialPolicy
	logger        *slog.Logger
}

func newProcessor(retrialPolicy RetrialPolicy, store Store, messageBroker MessageBroker, machineID string, logger *slog.Logger) *defaultRecordProcessor {
	return &defaultRecordProcessor{
		messageBroker: messageBroker,
		store:         store,
		time:          time.NewTimeProvider(),
		machineID:     machineID,
		retrialPolicy: retrialPolicy,
		logger:        logger,
	}
}

func (d defaultRecordProcessor) ProcessRecords() error {
	const op = "outbox.record_processor.ProcessRecords"

	defer func() {
		if cErr := d.store.ClearLocksByLockID(d.machineID); cErr != nil {
			errlog.ErrLog(richerror.New(op).WithWrapError(cErr).WithKind(richerror.KindUnexpected), d.logger)
		}
	}()

	if lErr := d.lockUnprocessedEntities(); lErr != nil {
		return richerror.New(op).WithWrapError(lErr).WithKind(richerror.KindUnexpected)
	}

	records, gErr := d.store.GetRecordsByLockID(d.machineID)
	if gErr != nil {
		return richerror.New(op).WithWrapError(gErr).WithKind(richerror.KindUnexpected)
	}

	if len(records) == 0 {
		return nil
	}

	return d.publishMessages(records)
}

func (d defaultRecordProcessor) publishMessages(records []Record) error {
	const op = "outbox.record_processor.publishMessages"

	for _, rec := range records {
		now := d.time.Now()
		rec.LastAttemptOn = &now
		rec.NumberOfAttempts++

		sErr := d.messageBroker.Send(rec.Message)
		if sErr != nil {
			rec.LockedOn = nil
			rec.LockID = nil
			errorMsg := sErr.Error()
			rec.Error = &errorMsg

			if d.retrialPolicy.MaxSendAttemptsEnabled && rec.NumberOfAttempts == d.retrialPolicy.MaxSendAttempts {
				rec.State = StateMaxAttemptsReached
			}

			if uErr := d.store.UpdateRecordByID(rec); uErr != nil {
				return richerror.New(op).WithMessage(fmt.Sprintf("could not update the record in the db: %s", uErr.Error())).
					WithWrapError(uErr).WithKind(richerror.KindUnexpected)
			}

			return richerror.New(op).WithWrapError(sErr).WithKind(richerror.KindUnexpected).
				WithMessage(fmt.Sprintf("an error occurred when trying to send the message to the broker: %s", sErr.Error()))
		}

		rec.State = StateDelivered
		rec.LockedOn = nil
		rec.LockID = nil
		rec.ProcessedOn = &now

		if uErr := d.store.UpdateRecordByID(rec); uErr != nil {
			return richerror.New(op).WithWrapError(uErr).WithKind(richerror.KindUnexpected).
				WithMessage(fmt.Sprintf("could not update the record in the db: %s", uErr.Error()))
		}
	}

	return nil
}

func (d defaultRecordProcessor) lockUnprocessedEntities() error {
	const op = "outbox.record_processor.lockUnprocessedEntities"

	lockTime := d.time.Now().UTC()

	if lErr := d.store.UpdateRecordLockByState(d.machineID, lockTime, StatePendingDelivery); lErr != nil {
		return richerror.New(op).WithWrapError(lErr).WithKind(richerror.KindUnexpected)
	}

	return nil
}
