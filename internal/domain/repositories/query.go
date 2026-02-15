package repositories

import "time"

// SortDirection represents ASC or DESC ordering.
type SortDirection string

const (
	SortAsc  SortDirection = "asc"
	SortDesc SortDirection = "desc"
)

// ChemistryLogQuery holds pagination, sorting, and filter parameters.
type ChemistryLogQuery struct {
	Page         int
	PageSize     int
	SortBy       string
	SortDir      SortDirection
	OutOfRange   bool
	DateFrom     *time.Time
	DateTo       *time.Time
}

// Defaults fills zero values with sensible defaults.
func (q *ChemistryLogQuery) Defaults() {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 {
		q.PageSize = 25
	}
	if q.PageSize > 100 {
		q.PageSize = 100
	}
	if q.SortBy == "" {
		q.SortBy = "tested_at"
	}
	if q.SortDir == "" {
		q.SortDir = SortDesc
	}
}

// Offset returns the SQL OFFSET for the current page.
func (q *ChemistryLogQuery) Offset() int {
	return (q.Page - 1) * q.PageSize
}

// PagedResult holds a page of results plus pagination metadata.
type PagedResult[T any] struct {
	Items      []T
	Page       int
	PageSize   int
	TotalItems int
	TotalPages int
}
