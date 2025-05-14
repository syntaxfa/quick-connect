package outbox

import (
	"log/slog"
	"os"
	"time"

	"github.com/syntaxfa/quick-connect/pkg/errlog"
)

type processor interface {
	ProcessRecords() error
}

type unlocker interface {
	UnlockExpiredMessages() error
}

type cleaner interface {
	RemoveExpiredMessages() error
}

type Dispatcher struct {
	recordProcessor processor
	recordUnlocker  unlocker
	recordCleaner   cleaner
	cfg             Config
	logger          *slog.Logger
}

func NewDispatcher(cfg Config, store Store, broker MessageBroker, machineID string, logger *slog.Logger) Dispatcher {
	return Dispatcher{
		recordProcessor: newProcessor(cfg.RetrialPolicy, store, broker, machineID, logger),
		recordUnlocker:  newRecordUnlocker(store, cfg.MaxLockTimeDuration),
		recordCleaner:   newRecordCleaner(store, cfg.MessagesRetentionDuration),
		cfg:             cfg,
		logger:          logger,
	}
}

func (d Dispatcher) Run(trap chan os.Signal) {
	go d.runRecordProcessor(trap)
	go d.runRecordUnlocker(trap)
	go d.runRecordCleaner(trap)

	<-trap

	d.logger.Info("stopping dispatcher")
}

func (d Dispatcher) runRecordProcessor(trap chan os.Signal) {
	ticker := time.NewTicker(d.cfg.ProcessInterval)

	for {
		d.logger.Info("record processor running!!!")

		pErr := d.recordProcessor.ProcessRecords()
		if pErr != nil {
			errlog.WithoutErr(pErr, d.logger)
		}

		d.logger.Info("record processing finished")

		select {
		case <-ticker.C:
			continue
		case <-trap:
			ticker.Stop()
			d.logger.Info("stopping record processor")

			return
		}
	}
}

func (d Dispatcher) runRecordUnlocker(trap chan os.Signal) {
	ticker := time.NewTicker(d.cfg.LockCheckerInterval)

	for {
		d.logger.Info("record unlocker running")

		if uErr := d.recordUnlocker.UnlockExpiredMessages(); uErr != nil {
			errlog.WithoutErr(uErr, d.logger)
		}

		d.logger.Info("record unlocker finished")

		select {
		case <-ticker.C:
			continue
		case <-trap:
			ticker.Stop()
			d.logger.Info("stopping record unlocker")

			return
		}
	}
}

func (d Dispatcher) runRecordCleaner(trap chan os.Signal) {
	ticker := time.NewTicker(d.cfg.CleanupWorkerInterval)

	for {
		d.logger.Info("record retention cleaner running")

		if rErr := d.recordCleaner.RemoveExpiredMessages(); rErr != nil {
			errlog.WithoutErr(rErr, d.logger)
		}

		d.logger.Info("record retention cleaner finished")

		select {
		case <-ticker.C:
			continue
		case <-trap:
			ticker.Stop()
			d.logger.Info("stopping record retention cleaner")

			return
		}
	}
}
