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

	// --- Refactored Part ---
	// Build WHERE clause using a helper function
	paramIndex := 1
	whereClause, whereArgs, nextParamIndex := buildConditions(parameters.Filters, paramIndex)

	args = append(args, whereArgs...)
	paramIndex = nextParamIndex

	if whereClause != "" {
		query += " WHERE " + whereClause
		countQuery += " WHERE " + whereClause + ";" // Add semicolon for count query
	}
	// --- End Refactored Part ---

	// User "ID" as the default column if sortColumn is empty
	if parameters.SortColumn == "" {
		parameters.SortColumn = DefaultSortColumn
	}

	// Add sorting and pagination
	query += fmt.Sprintf(" ORDER BY %s %s", parameters.SortColumn, orderDirection(parameters.Descending))
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)

	// Add pagination arguments
	args = append(args, parameters.Limit, parameters.Offset)

	return query, countQuery, args
}

// buildConditions builds the WHERE clause, arguments, and returns the next parameter index.
func buildConditions(filters map[paginate.FilterParameter]paginate.Filter, startIndex int) (string, []interface{}, int) {
	if len(filters) == 0 {
		return "", nil, startIndex
	}

	conditions := make([]string, 0, len(filters))
	var args []interface{}
	paramIndex := startIndex

	for p, f := range filters {
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
			placeholders := make([]string, 0, len(f.Values))
			for _, value := range f.Values {
				placeholders = append(placeholders, fmt.Sprintf("$%d", paramIndex))
				args = append(args, value)
				paramIndex++
			}
			conditions = append(conditions, fmt.Sprintf("%s IN (%s)", p, strings.Join(placeholders, ", ")))
		case paginate.FilterOperationNotIn:
			placeholders := make([]string, 0, len(f.Values))
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

	return strings.Join(conditions, " AND "), args, paramIndex
}

func orderDirection(descending bool) string {
	if descending {
		return "DESC"
	}

	return "ASC"
}
