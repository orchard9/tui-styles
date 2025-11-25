package domain

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// PaginationType defines the pagination strategy
type PaginationType string

const (
	// PaginationTypeOffset uses LIMIT/OFFSET (simpler, less performant at scale)
	PaginationTypeOffset PaginationType = "offset"

	// PaginationTypeCursor uses cursor-based pagination (performant, consistent)
	PaginationTypeCursor PaginationType = "cursor"
)

// SortDirection defines sort order
type SortDirection string

const (
	SortDirectionAsc  SortDirection = "asc"
	SortDirectionDesc SortDirection = "desc"
)

// SortField defines a field to sort by
type SortField struct {
	Field     string        // Field name (e.g., "created_at", "email")
	Direction SortDirection // Sort direction
}

// Filter defines a single filter condition
type Filter struct {
	Field    string      // Field name
	Operator string      // Operator: eq, ne, gt, gte, lt, lte, like, in
	Value    interface{} // Filter value
}

// ListOptions contains all query parameters for list operations
type ListOptions struct {
	// Pagination strategy
	Type PaginationType

	// Limit: Maximum number of results (1-100)
	Limit int

	// Offset: Number of records to skip (offset pagination only)
	Offset int

	// Cursor: Opaque token for cursor pagination
	Cursor string

	// SortBy: Fields to sort by (in order)
	SortBy []SortField

	// Filters: Conditions to apply
	Filters []Filter

	// IncludeTotal: Whether to include total count
	IncludeTotal bool
}

// Validate checks if ListOptions are valid
func (opts *ListOptions) Validate() error {
	// Validate limit
	if opts.Limit <= 0 {
		opts.Limit = 20 // Default
	}
	if opts.Limit > 100 {
		return ErrInvalidInput("limit cannot exceed 100")
	}

	// Validate offset (if using offset pagination)
	if opts.Type == PaginationTypeOffset && opts.Offset < 0 {
		return ErrInvalidInput("offset cannot be negative")
	}

	// Validate cursor (if using cursor pagination)
	if opts.Type == PaginationTypeCursor && opts.Cursor != "" {
		if _, err := DecodeCursor(opts.Cursor); err != nil {
			return ErrInvalidInput(fmt.Sprintf("invalid cursor: %v", err))
		}
	}

	// Validate sort fields
	for _, sort := range opts.SortBy {
		if sort.Field == "" {
			return ErrInvalidInput("sort field cannot be empty")
		}
		if sort.Direction != SortDirectionAsc && sort.Direction != SortDirectionDesc {
			return ErrInvalidInput(fmt.Sprintf("invalid sort direction: %s", sort.Direction))
		}
	}

	// Set default sort if not specified
	if len(opts.SortBy) == 0 {
		opts.SortBy = []SortField{
			{Field: "created_at", Direction: SortDirectionDesc},
		}
	}

	return nil
}

// ListResult contains paginated results with metadata
type ListResult struct {
	// Total: Total count of all results (ignoring pagination)
	// Only populated if IncludeTotal=true
	// Set to -1 if not calculated
	Total int

	// HasMore: Whether there are more results after this page
	HasMore bool

	// NextCursor: Cursor for next page (cursor pagination)
	NextCursor string

	// PrevCursor: Cursor for previous page (cursor pagination)
	PrevCursor string

	// Limit: Limit used for this query (echoed back)
	Limit int

	// Offset: Offset used for this query (offset pagination, echoed back)
	Offset int
}

// Cursor represents a cursor for pagination
type Cursor struct {
	ID         string                 `json:"id"`
	CreatedAt  int64                  `json:"created_at"`
	SortValues map[string]interface{} `json:"sort_values,omitempty"`
}

// EncodeCursor encodes a cursor to base64 string
func EncodeCursor(c Cursor) string {
	data, _ := json.Marshal(c)
	return base64.StdEncoding.EncodeToString(data)
}

// DecodeCursor decodes a base64 cursor string
func DecodeCursor(encoded string) (*Cursor, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor encoding: %w", err)
	}

	var c Cursor
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("invalid cursor format: %w", err)
	}

	return &c, nil
}
