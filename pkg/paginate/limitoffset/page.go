package paginate

var (
	DefaultMinPageSize uint64 = 10
	DefaultMaxPageSize uint64 = 100
)

// FilterParameter defines a key for filtering paginated results.
// Custom filter parameters can be defined as needed by the application in service layer.
// Example value include:
//   - "status": To filter results by status (e.g, "active", "inactive").
//   - "user_id": To filter results by a specific user ID.
//   - "created_at": To filter results based on creation date.
type FilterParameter string

// Paginated struct represents a paginated response that will be passed to the repository layer
// to retrieve data based on pagination settings, filters, sorting, etc.
type Paginated struct {
	Page       uint64                     `json:"page"`
	PerPage    uint64                     `json:"per_page"`
	Total      uint64                     `json:"total"`
	Filters    map[FilterParameter]Filter `json:"filters"`
	SortColumn string                     `json:"sort_column"`
	Descending bool                       `json:"descending"`
}

/*
Example of how paginated struct can be used in a database implementation
type PaginationSupportDB interface {
    GetPaginated(ctx context.Context, p Paginated) (res []Result, total int64, err error)
}

This interface defines how pagination and filtering data can be fetched from the database.
*/

type Filter struct {
	Operation FilterOperation
	Values    []interface{}
}

type FilterOperation int

const (
	FilterOperationEqual FilterOperation = iota + 1
	FilterOperationNotEqual
	FilterOperationGreater
	FilterOperationGreaterEqual
	FilterOperationLess
	FilterOperationLessEqual
	FilterOperationIn
	FilterOperationNotIn
	FilterOperationBetween
)

// RequestBase the base structure for pagination requests from the client.
type RequestBase struct {
	CurrentPage uint64                     `json:"current_page"`
	PageSize    uint64                     `json:"page_size"`
	Filters     map[FilterParameter]Filter `json:"filters"`
	SortColumn  string                     `json:"sort_column"`
	Descending  bool                       `json:"descending"`
}

// ResponseBase the base structure for paginated responses sent to client.
type ResponseBase struct {
	CurrentPage  uint64 `json:"current_page"`
	PageSize     uint64 `json:"page_size"`
	TotalNumbers uint64 `json:"total_numbers"`
	TotalPage    uint64 `json:"total_page"`
}

// BasicValidation just ensure that the pagination request is well-formed.
// more complex validation should be dome in service layers based on the application's requirements.
func (r *RequestBase) BasicValidation() error {
	if r.CurrentPage < 1 {
		r.CurrentPage = 1
	}

	if r.PageSize < DefaultMinPageSize {
		r.PageSize = DefaultMinPageSize
	}

	if r.PageSize > DefaultMaxPageSize {
		r.PageSize = DefaultMaxPageSize
	}

	// TODO: Add filters validation to ensure that the filters are well-formed

	return nil
}
