package postgres

import "database/sql"

type nullableFields struct {
	Bucket    sql.NullString
	DeletedAt sql.NullTime
}
