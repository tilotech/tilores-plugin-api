package api

// Entity represents a real world object
type Entity struct {
	ID         string
	Records    []*Record
	Edges      Edges
	Duplicates Duplicates
}

// Edges represents a connection between two Records
//
// e.g. "recordID:anotherRecordID:STATIC" or "recordID:anotherRecordID:RULEID"
type Edges []string

// Duplicates represents all record duplicates within the entity
//
// e.g. {"record1ID":["duplicateOfRecord1ID"],"record2ID":["firstDuplicateOfRecord2ID", "secondDuplicateOfRecord2ID"]}
type Duplicates map[string][]string
