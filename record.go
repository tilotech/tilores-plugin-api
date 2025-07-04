package api

import (
	"time"
)

// Record represents a part of an Entity and the corresponding predicates
//
// Each Record must have a unique ID.
type Record struct {
	ID   string         `json:"id"`
	Data map[string]any `json:"data"`
	Meta *RecordMeta    `json:"meta"`
}

// RecordMeta stores additional information about the record.
type RecordMeta struct {
	SubmitTimestamp   *time.Time          `json:"submitTimestamp"`
	AssembleTimestamp *time.Time          `json:"assembleTimestamp"`
	Version           int                 `json:"version"`
	ConsistencyIndex  map[string][]string `json:"consistencyIndex"`
}

// IDWithVersion returns the records ID and its version in the format <id>:<version>.
func (r *Record) IDWithVersion() string {
	v := 0
	if r.Meta != nil {
		v = r.Meta.Version
	}
	return NewRecordID(r.ID, v)
}
