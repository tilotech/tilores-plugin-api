package api

import (
	"fmt"
	"strconv"
	"strings"
)

// Entity represents a real world object
type Entity struct {
	ID          string     `json:"id"`
	Records     []*Record  `json:"records"`
	Edges       Edges      `json:"edges"`
	Duplicates  Duplicates `json:"duplicates"`
	Hits        Hits       `json:"hits"`
	Consistency float64    `json:"consistency"`
	Score       float64    `json:"score"`
	HitScore    float64    `json:"hitScore"`
}

// Edges represents a connection between two Records
//
// Edges are typically in the form <id>:<version>:<id>:<version>:<rule>
// e.g. "recordID:0:anotherRecordID:2:STATIC" or "recordID:0:anotherRecordID:2:RULEID"
//
// Edges displayed to the user omit the version information, resulting in the
// following form: <id>:<id>:<rule>
type Edges []string

// Duplicates represents all record duplicates within the entity
//
// Duplicates are typically in the form <group>:<id>:<version> with <group>
// allowed to be empty.
// e.g. {":record1ID:2":["duplicateOfRecord1ID"],"record2ID":["firstDuplicateOfRecord2ID", "secondDuplicateOfRecord2ID"]}
//
// Duplicates displayed to the user omit the version information and empty
// groups, resulting in one of the following forms:
// <id> or <group>:<id>
type Duplicates map[string][]string

// Hits lists all matched rules per matched record id
//
// Example (in JSON):
//
//	{
//	  "550e8400-e29b-11d4-a716-446655440000": ["RULE-1", "RULE-2"],
//	  "6ba7b810-9dad-11d1-80b4-00c04fd430c8:2": ["RULE-2"],
//	}
//
// The maps key is a record ID as specified in NewRecordID.
type Hits map[string][]string

// IDs returns the record ids of the hits
func (h Hits) IDs() []string {
	ids := make([]string, len(h))
	i := 0
	for k := range h {
		rid, _ := ParseRecordID(k)
		ids[i] = rid
		i++
	}
	return ids
}

// ParseEdge parses an edge string into its components.
//
// Return values are: id1, id2 and rule
// Both id1 and id2 are record IDs as defined by NewRecordID.
//
// The edge string must be in the format:
// <id>:<version>:<id>:<version>:<rule><score>,
// <id>:<version>:<id>:<version>:<rule> or <id>:<id>:<rule>
// If the version information is not provided, then version 0 is assumed.
// If the score is not provided, then score 100 is assumed.
// The behavior for other formats is undefined.
func ParseEdge(edge string) (string, string, string, uint8) {
	parts := strings.SplitN(edge, ":", 6)
	var id1, id2, rule string
	var v1, v2 int
	var score uint8
	switch len(parts) {
	case 3:
		id1, v1 = ParseRecordID(parts[0])
		id2, v2 = ParseRecordID(parts[1])
		rule = parts[2]
		score = 100
	case 5:
		id1, v1 = ParseRecordID(fmt.Sprintf("%v:%v", parts[0], parts[1]))
		id2, v2 = ParseRecordID(fmt.Sprintf("%v:%v", parts[2], parts[3]))
		rule = parts[4]
		score = 100
	default:
		id1, v1 = ParseRecordID(fmt.Sprintf("%v:%v", parts[0], parts[1]))
		id2, v2 = ParseRecordID(fmt.Sprintf("%v:%v", parts[2], parts[3]))
		rule = parts[4]
		if s, err := strconv.ParseUint(parts[5], 10, 8); err == nil {
			score = uint8(s)
		}
	}
	return NewRecordID(id1, v1), NewRecordID(id2, v2), rule, score
}

// NewEdge returns a new edge string with the provided IDs and versions.
//
// id1 and id2 must be in the format: <id>:<version> or <id>
// If id1 or id2 is in the format <id>, then the version 0 is assumed.
//
// score must be a value between 0 and 100 (both including).
//
// The resulting string will be in the format: <id>:<version>:<id>:<version>:<rule>:<score>
func NewEdge(id1, id2, rule string, score uint8) string {
	id1, v1 := ParseRecordID(id1)
	id2, v2 := ParseRecordID(id2)
	return NewEdgeWithVersions(id1, v1, id2, v2, rule, score)
}

// NewEdgeWithVersions returns a new edge string with the provided IDs and versions.
//
// id1 and id2 must be the plain record ID without a version.
//
// score must be a value between 0 and 100 (both including).
//
// The resulting string will be in the format: <id>:<version>:<id>:<version>:<rule>:<score>
func NewEdgeWithVersions(id1 string, v1 int, id2 string, v2 int, rule string, score uint8) string {
	return fmt.Sprintf("%v:%v:%v:%v:%v:%v", id1, v1, id2, v2, rule, score)
}

// ParseDuplicateKey parses the key (original) of a duplicate into its components.
//
// Return values are: id, group
// id is a record ID as defined by NewRecordID.
//
// The duplicate key must be in the format:
// <group>:<id>:<version> or <group>:<id> or <id>
// If the <group> is not present, then an empty string is assumed.
// If the version is not present, then version 0 is assumed.
//
// The group can be an empty string.
// The behavior for other formats is undefined, especially the format
// <id>:<version> is not defined!
func ParseDuplicateKey(key string) (string, string) {
	parts := strings.SplitN(key, ":", 3)
	var group, id string
	var v int
	switch len(parts) {
	case 1:
		id = parts[0]
	case 2:
		group = parts[0]
		id = parts[1]
	default:
		group = parts[0]
		id, v = ParseRecordID(fmt.Sprintf("%v:%v", parts[1], parts[2]))
	}
	return NewRecordID(id, v), group
}

// NewDuplicateKey returns a new duplicate key string with the provided id and group.
//
// id must be in the format: <id>:<version> or <id>
// If id is in the format <id>, then version 0 is assumed.
// The group can be an empty string.
//
// The resulting string will be in the format: <group>:<id>:<version>
func NewDuplicateKey(id, group string) string {
	rid, v := ParseRecordID(id)
	return fmt.Sprintf("%v:%v:%v", group, rid, v)
}

// ParseRecordID parses the record id into its components.
//
// Return values are: id, version
//
// The recordID must be in the format: <id>:<version> or <id>
// If only <id> is used, then the version is assumed to be 0.
// The behavior for other formats is undefined.
func ParseRecordID(recordID string) (string, int) {
	parts := strings.SplitN(recordID, ":", 2)
	if len(parts) == 1 {
		return parts[0], 0
	}
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

// NewRecordID returns a new record id string with the provided id and
// optional version.
//
// The resulting string will be in the format: <id>:<version> or <id>
// If the provided version is equal to 0, then the output will not include the
// version.
func NewRecordID(id string, version int) string {
	if version == 0 {
		return id
	}
	return fmt.Sprintf("%v:%v", id, version)
}

// NewRecordIDWithVersion returns a new record id string with the provided id
// and version.
//
// The resulting string will be in the format: <id>:<version>
// This is different from the behavior of NewRecordID that omits the version
// if it is equal to 0.
//
// Most cases should use NewRecordID.
func NewRecordIDWithVersion(id string, version int) string {
	return fmt.Sprintf("%v:%v", id, version)
}
