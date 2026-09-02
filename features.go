package api

// Features defines the features that are active and must be calculated
// or collected and returned during a query.
//
// A feature is only considered inactive in case it was explicitly set to false.
// The feature is considered active in case of nil, missing value and explicit true.
type Features struct {
	EntityConsistency *bool `json:"entityConsistency"`
	// EntityConsistencyValues is the one exception to the rule above: a NIL
	// flag counts as INACTIVE, not active.
	//
	// The others must read nil as active because a client old enough to send no
	// feature flags at all would otherwise lose data it was asking for. That
	// cannot happen here: the field only exists on a server new enough to
	// report it, and such a server always sets the flag explicitly. Reading nil
	// as active would instead make every older client pay to resolve every
	// record through the consistency rule set for values it cannot return.
	//
	// A hand-written client that wants the values must therefore set this
	// explicitly; leaving it unset yields an empty list rather than an error.
	EntityConsistencyValues *bool `json:"entityConsistencyValues"`
	EntityDuplicates        *bool `json:"entityDuplicates"`
	EntityEdges             *bool `json:"entityEdges"`
	EntityHits              *bool `json:"entityHits"`
	EntityHitScore          *bool `json:"entityHitScore"`
	EntityRecords           *bool `json:"entityRecords"`
	EntityScore             *bool `json:"entityScore"`
}
