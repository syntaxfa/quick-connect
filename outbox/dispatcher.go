package outbox

import (
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/logger"
	"os"
	"time"
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
}

func NewDispatcher(cfg Config, store Store, broker MessageBroker, machineID string) Dispatcher {
	return Dispatcher{
		recordProcessor: newProcessor(cfg.RetrialPolicy, store, broker, machineID),
		recordUnlocker:  newRecordUnlocker(store, cfg.MaxLockTimeDuration),
		recordCleaner:   newRecordCleaner(store, cfg.MessagesRetentionDuration),
		cfg:             cfg,
	}
}

func (d Dispatcher) Run(trap chan os.Signal) {
	go d.runRecordProcessor(trap)
	go d.runRecordUnlocker(trap)
	go d.runRecordCleaner(trap)

	<-trap

	logger.L().Info("stopping dispatcher")
}

func (d Dispatcher) runRecordProcessor(trap chan os.Signal) {
	ticker := time.NewTicker(d.cfg.ProcessInterval)

	for {
		logger.L().Info("record processor running!!!")

		pErr := d.recordProcessor.ProcessRecords()
		if pErr != nil {
			errlog.ErrLog(pErr)
		}

		logger.L().Info("record processing finished")

		select {
		case <-ticker.C:
			continue
		case <-trap:
			ticker.Stop()
			logger.L().Info("stopping record processor")

			return
		}
	}
}

func (d Dispatcher) runRecordUnlocker(trap chan os.Signal) {
	ticker := time.NewTicker(d.cfg.LockCheckerInterval)

	for {
		logger.L().Info("record unlocker running")

		if uErr := d.recordUnlocker.UnlockExpiredMessages(); uErr != nil {
			errlog.ErrLog(uErr)
		}

		logger.L().Info("record unlocker finished")

		select {
		case <-ticker.C:
			continue
		case <-trap:
			ticker.Stop()
			logger.L().Info("stopping record unlocker")

			return
		}
	}
}

func (d Dispatcher) runRecordCleaner(trap chan os.Signal) {
	ticker := time.NewTicker(d.cfg.CleanupWorkerInterval)

	for {
		logger.L().Info("record retention cleaner running")

		if rErr := d.recordCleaner.RemoveExpiredMessages(); rErr != nil {
			errlog.ErrLog(rErr)
		}

		logger.L().Info("record retention cleaner finished")

		select {
		case <-ticker.C:
			continue
		case <-trap:
			ticker.Stop()
			logger.L().Info("stopping record retention cleaner")

			return
		}
	}
}
