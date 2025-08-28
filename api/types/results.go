package types

// FindResult is a generic container for find operations that return a list of
// items and a total count.
type FindResult struct {
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
}

// NewFindResult creates a new FindResult.
func NewFindResult(data interface{}, total int64) *FindResult {
	return &FindResult{
		Data:  data,
		Total: total,
	}
}
