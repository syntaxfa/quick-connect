package postgres

type statementKey uint

const (
	StatementClearLocksWithDurationBeforeDate statementKey = iota + 1
	StatementUpdateRecordByID
	StatementUpdateRecordLockByState
	StatementClearLocksByLockID
	StatementGetRecordsByLockID
	StatementAddRecordTX
	StatementRemoveRecordsBeforeDatetime
)
