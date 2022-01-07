package api

// Record represents a part of an Entity and the corresponding predicates
//
// Each Record must have a unique ID.
type Record struct {
	ID   string                 `json:"id"`
	Data map[string]interface{} `json:"data"`
}
