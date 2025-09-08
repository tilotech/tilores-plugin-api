package api

// Features defines the features that are active and must be calculated
// or collected and returned during a query.
//
// A feature is only considered inactive in case it was explicitly set to false.
// The feature is considered active in case of nil, missing value and explicit true.
type Features struct {
	EntityConsistency *bool `json:"entityConsistency"`
}
