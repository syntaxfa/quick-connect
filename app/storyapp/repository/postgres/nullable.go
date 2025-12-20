package postgres

import "database/sql"

type nullableFields struct {
	Title    sql.NullString
	Caption  sql.NullString
	LinkURL  sql.NullString
	LinkText sql.NullString
}
