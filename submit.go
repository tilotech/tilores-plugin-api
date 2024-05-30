package api

// SubmissionPreview a preview of how the data could potentially
// look like if/when submitted.
//
// Entities a list of entities where the provided records reside.
//
// NewRecords a list of record IDs that where newly added to resulting Entities.
//
// UpdatedRecords a list of record IDs that where updated in resulting Entities.
//
// IgnoredRecords a list of record IDs that where ignored and not/will not be
// ingested. This is only relevant in case record updates are not enabled.
type SubmissionPreview struct {
	Entities       []*Entity `json:"entities"`
	NewRecords     []string  `json:"newRecords"`
	UpdatedRecords []string  `json:"updatedRecords"`
	IgnoredRecords []string  `json:"ignoredRecords"`
}
