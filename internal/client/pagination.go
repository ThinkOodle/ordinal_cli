package client

import (
	"fmt"
	"net/url"
)

// CursorPaginationParams holds cursor-based pagination parameters.
// Ordinal uses {limit, cursor} with {nextCursor, hasMore} in the response.
type CursorPaginationParams struct {
	Limit  int
	Cursor string
}

// ApplyToQuery adds pagination parameters to URL query values.
func (p CursorPaginationParams) ApplyToQuery(q url.Values) {
	if p.Limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", p.Limit))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
}
