package api

import "time"

// Record represents a part of an Entity and the corresponding predicates
//
// Each Record must have a unique ID.
type Record struct {
	ID   string                 `json:"id"`
	Data map[string]interface{} `json:"data"`
	Meta *RecordMeta            `json:"meta"`
}

// RecordMeta stores additional information about the record.
type RecordMeta struct {
	SubmitTimestamp   *time.Time `json:"submitTimestamp"`
	AssembleTimestamp *time.Time `json:"assembleTimestamp"`
	Version           int        `json:"version"`
}
