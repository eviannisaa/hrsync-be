package model

// ListParams holds query parameters for list endpoints
type ListParams struct {
	Page    int
	Limit   int
	Search  string
	Email   string
	Role    string
	SortBy  string
	SortDir string // "asc" or "desc"
}
