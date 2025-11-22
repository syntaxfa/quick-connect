package cursorbased

import (
	"fmt"

	"github.com/oklog/ulid/v2"
	"github.com/syntaxfa/quick-connect/types"
)

const (
	DefaultLimit int = 20
	MaxLimit     int = 100
)

// Request represents the structure for cursor-based pagination requests.
type Request struct {
	// Cursor is unique identifier (ULID) of the last item from the previous page.
	// If empty, it includes the first page
	Cursor types.ID `json:"cursor"`
	// Limit specifies how many items to fetch.
	Limit int `json:"limit"`
}

// Response represents the structure for cursor-based pagination responses.
type Response struct {
	// Next Counter is the ID of the last item in the current list.
	// The client should user this as the `cursor` for the next request.
	NextCursor types.ID `json:"next_cursor"`
	// HasMore indicates if there are mote item available to fetch.
	HasMore bool `json:"has_more"`
}

// BasicValidation ensures the request parameters are within valid ranges.
func (r *Request) BasicValidation() error {
	if r.Limit < 1 {
		r.Limit = DefaultLimit
	}

	if r.Limit > MaxLimit {
		r.Limit = MaxLimit
	}

	if r.Cursor != "" {
		if _, pErr := ulid.Parse(string(r.Cursor)); pErr != nil {
			return fmt.Errorf("cursor %s is not a valid ulid", r.Cursor)
		}
	}

	// Note: Empty cursor is valid (means first page)
	return nil
}
