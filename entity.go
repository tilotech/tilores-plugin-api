package api

import (
	"fmt"
	"strconv"
	"strings"
)

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
		ids[i] = strings.SplitN(k, ":", 2)[0]
		i++
	}
	return ids
}

// ParseEdge parses an edge string into its components.
//
// Return values are: id1, id2 and rule
// Both id1 and id2 will contain the version string, e.g. foo:1
//
// The edge string must be in the format: <id>:<version>:<id>:<version>:<rule>
// The behavior for other formats is undefined.
func ParseEdge(edge string) (string, string, string) {
	parts := strings.SplitN(edge, ":", 5)
	return fmt.Sprintf("%v:%v", parts[0], parts[1]), fmt.Sprintf("%v:%v", parts[2], parts[3]), parts[4]
}

// NewEdge returns a new edge string with the provided IDs and versions.
//
// id1 and id2 must be in the format: <id>:<version>
//
// The resulting string will be in the format: <id>:<version>:<id>:<version>:<rule>
func NewEdge(id1, id2, rule string) string {
	return fmt.Sprintf("%v:%v:%v", id1, id2, rule)
}

// NewEdgeWithVersions returns a new edge string with the provided IDs and versions.
//
// id1 and id2 must be the plain record ID without a version.
//
// The resulting string will be in the format: <id>:<version>:<id>:<version>:<rule>
func NewEdgeWithVersions(id1 string, v1 int, id2 string, v2 int, rule string) string {
	return fmt.Sprintf("%v:%v:%v:%v:%v", id1, v1, id2, v2, rule)
}

// ParseDuplicateKey parses the key (original) of a duplicate into its components.
//
// Return values are: id, group
// Id will contain the version string, e.g. foo:1
//
// The duplicate key must be in the format: <group>:<id>:<version>
// The group can be an empty string.
// The behavior for other formats is undefined.
func ParseDuplicateKey(key string) (string, string) {
	parts := strings.SplitN(key, ":", 3)
	return fmt.Sprintf("%v:%v", parts[1], parts[2]), parts[0]
}

// NewDuplicateKey returns a new duplicate key string with the provided id and group.
//
// id must be in the format: <id>:<version>
// The group can be an empty string.
//
// The resulting string will be in the format: <group>:<id>:<version>
func NewDuplicateKey(id, group string) string {
	return fmt.Sprintf("%v:%v", group, id)
}

// ParseRecordID parses the record id into its components.
//
// Return values are: id, version
//
// The recordID must be in the format: <id>:<version>
// The behavior for other formats is undefined.
func ParseRecordID(recordID string) (string, int) {
	parts := strings.SplitN(recordID, ":", 2)
	version, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}
	return parts[0], version
}

// ParseRecordIDWithOptionalVersion parses the record id into its components.
//
// Return values are: id, version
//
// The recordID must be in the format: <id>:<version> or <id>
// The behavior for other formats is undefined.
func ParseRecordIDWithOptionalVersion(recordID string) (string, *int) {
	parts := strings.SplitN(recordID, ":", 2)
	if len(parts) == 1 {
		return parts[0], nil
	}
	version, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}
	return parts[0], &version
}

// NewRecordID returns a new record id string with the provided id and version.
//
// The resulting string will be in the format: <id>:<version>
func NewRecordID(id string, version int) string {
	return fmt.Sprintf("%v:%v", id, version)
}
