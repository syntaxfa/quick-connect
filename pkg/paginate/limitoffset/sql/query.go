package pagesql

import (
	"fmt"
	"strings"

	paginate "github.com/syntaxfa/quick-connect/pkg/paginate/limitoffset"
)

var DefaultSortColumn = "id"

type Parameters struct {
	Table      string
	Fields     []string
	Filters    map[paginate.FilterParameter]paginate.Filter
	SortColumn string
	Descending bool
	Limit      uint64
	Offset     uint64
}

// WriteQuery generates a SQL query for paginated results and a count query based on the provided filters, sorting, and pagination parameters.
// It constructs a SELECT query with conditions for filtering, ordering, pagination (LIMIT and OFFSET) and a COUNT query for the total number of records.
func WriteQuery(parameters Parameters) (query, countQuery string, args []interface{}) {
	selectFields := "*"
	if len(parameters.Fields) > 0 {
		selectFields = strings.Join(parameters.Fields, ", ")
	}

	// Base query for pagination
	query = fmt.Sprintf("SELECT %s FROM %s", selectFields, parameters.Table)

	// Base query for total count
	countQuery = fmt.Sprintf("SELECT COUNT(*) FROM %s", parameters.Table)

	// Add filters to the query
	if len(parameters.Filters) == 0 {
		// User DefaultSortColumn if sortColumn is empty
		if parameters.SortColumn == "" {
			sortColumn := DefaultSortColumn

			orderClause := fmt.Sprintf(" ORDER BY %s %s", sortColumn, orderDirection(parameters.Descending))

			// Set the pagination arguments
			args = []interface{}{parameters.Limit, parameters.Offset}

			return query + orderClause, countQuery, args
		}
	}

	query += " WHERE "
	conditions := make([]string, 0)
	paramIndex := len(args) + 1 // Start parameter numbering after pagination arguments

	for p, f := range parameters.Filters {
		switch f.Operation {
		case paginate.FilterOperationEqual:
			conditions = append(conditions, fmt.Sprintf("%s = $%d", p, paramIndex))
			args = append(args, f.Values[0])
			paramIndex++
		case paginate.FilterOperationNotEqual:
			conditions = append(conditions, fmt.Sprintf("%s != $%d", p, paramIndex))
			args = append(args, f.Values[0])
			paramIndex++
		case paginate.FilterOperationGreater:
			conditions = append(conditions, fmt.Sprintf("%s > $%d", p, paramIndex))
			args = append(args, f.Values[0])
			paramIndex++
		case paginate.FilterOperationGreaterEqual:
			conditions = append(conditions, fmt.Sprintf("%s >= $%d", p, paramIndex))
			args = append(args, f.Values[0])
			paramIndex++
		case paginate.FilterOperationLess:
			conditions = append(conditions, fmt.Sprintf("%s < $%d", p, paramIndex))
			args = append(args, f.Values[0])
			paramIndex++
		case paginate.FilterOperationLessEqual:
			conditions = append(conditions, fmt.Sprintf("%s <= $%d", p, paramIndex))
			args = append(args, f.Values[0])
			paramIndex++
		case paginate.FilterOperationIn:
			placeholders := make([]string, 0)
			for _, value := range f.Values {
				placeholders = append(placeholders, fmt.Sprintf("$%d", paramIndex))
				args = append(args, value)
				paramIndex++
			}
			conditions = append(conditions, fmt.Sprintf("%s IN (%s)", p, strings.Join(placeholders, ", ")))
		case paginate.FilterOperationNotIn:
			placeholders := make([]string, 0)
			for _, value := range f.Values {
				placeholders = append(placeholders, fmt.Sprintf("$%d", paramIndex))
				args = append(args, value)
				paramIndex++
			}
			conditions = append(conditions, fmt.Sprintf("%s NOT IN (%s)", p, strings.Join(placeholders, ", ")))
		case paginate.FilterOperationBetween:
			conditions = append(conditions, fmt.Sprintf("%s BETWEEN $%d AND $%d", p, paramIndex, paramIndex+1))
			args = append(args, f.Values[0], f.Values[1])
			paramIndex += 2
		}
	}

	// Add the condition to both the query and countQuery
	query += strings.Join(conditions, " AND ")

	// User "ID" as the default column if sortColumn is empty
	if parameters.SortColumn == "" {
		parameters.SortColumn = DefaultSortColumn
	}

	// Add sorting and pagination
	query += fmt.Sprintf(" ORDER BY %s %s", parameters.SortColumn, orderDirection(parameters.Descending))
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
	countQuery += fmt.Sprintf(" WHERE %s;", strings.Join(conditions, " AND "))

	// Add pagination arguments
	args = append(args, parameters.Limit, parameters.Offset)

	return query, countQuery, args
}

func orderDirection(descending bool) string {
	if descending {
		return "DESC"
	}

	return "ASC"
}
