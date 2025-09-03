package types

// FindResult is a generic container for find operations that return a list of
// items and a total count.
type FindResult[T any] struct {
	Data  []T   `json:"data"`
	Total int64 `json:"total"`
}

// NewFindResult creates a new FindResult.
func NewFindResult[T any](data []T, total int64) *FindResult[T] {
	return &FindResult[T]{
		Data:  data,
		Total: total,
	}
}
