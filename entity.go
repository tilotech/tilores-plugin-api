package api

// Entity represents a real world object
type Entity struct {
	ID      string
	Records []*Record
	Edges   Edges
}

// Edges represents a connection between two Records
//
// e.g. "recordID:anotherRecordID:STATIC" or "recordID:anotherRecordID:RULEID"
type Edges []string
