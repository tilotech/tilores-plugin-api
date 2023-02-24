package api

// Entity represents a real world object
type Entity struct {
	ID         string     `json:"id"`
	Records    []*Record  `json:"records"`
	Edges      Edges      `json:"edges"`
	Duplicates Duplicates `json:"duplicates"`
	Hits       Hits       `json:"hits"`
}

// Edges represents a connection between two Records
//
// e.g. "recordID:anotherRecordID:STATIC" or "recordID:anotherRecordID:RULEID"
type Edges []string

// Duplicates represents all record duplicates within the entity
//
// e.g. {"record1ID":["duplicateOfRecord1ID"],"record2ID":["firstDuplicateOfRecord2ID", "secondDuplicateOfRecord2ID"]}
type Duplicates map[string][]string

// Hits lists all matched rules per matched record id
//
// Example (in JSON):
//
//	{
//	  "550e8400-e29b-11d4-a716-446655440000": ["RULE-1", "RULE-2"],
//	  "6ba7b810-9dad-11d1-80b4-00c04fd430c8": ["RULE-2"],
//	}
type Hits map[string][]string

// IDs returns the record ids of the hits
func (h Hits) IDs() []string {
	ids := make([]string, len(h))
	i := 0
	for k := range h {
		ids[i] = k
		i++
	}
	return ids
}
