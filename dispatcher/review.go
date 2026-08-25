package dispatcher

import (
	"time"
)

// ReviewCaseStatus describes where a review case currently is in its lifecycle.
type ReviewCaseStatus string

const (
	// ReviewCaseOpen marks a case that nobody is working on.
	ReviewCaseOpen ReviewCaseStatus = "OPEN"
	// ReviewCaseClaimed marks a case that an actor is currently reviewing and
	// for which the affected entity is locked.
	ReviewCaseClaimed ReviewCaseStatus = "CLAIMED"
)

// ReviewCasesInput selects the review cases to return.
//
// Cases are returned in the order in which they were created, oldest first.
type ReviewCasesInput struct {
	Status *ReviewCaseStatus `json:"status"` // defaults to OPEN
	Cursor *string           `json:"cursor"` // opaque, empty starts at the beginning
	Limit  *int              `json:"limit"`
}

// ReviewCasesOutput is the result of a ReviewCases call.
type ReviewCasesOutput struct {
	Cases      []*ReviewCase `json:"cases"`
	NextCursor *string       `json:"nextCursor"` // nil means end of the queue
	OpenCount  int           `json:"openCount"`
}

// ReviewCase describes an assembly that joined one or more entities on evidence
// that was too weak to be trusted.
type ReviewCase struct {
	ID         string            `json:"id"`
	Status     ReviewCaseStatus  `json:"status"`
	CreatedAt  time.Time         `json:"createdAt"`
	RecordIDs  []string          `json:"recordIDs"` // the submitted records, with version
	EntityID   string            `json:"entityID"`  // the surviving entity at detection time
	Candidates []ReviewCandidate `json:"candidates"`
	ClaimedBy  *ReviewActor      `json:"claimedBy"`
	ClaimedAt  *time.Time        `json:"claimedAt"`
}

// ReviewCandidate is a single entity that the submitted records joined too
// weakly.
type ReviewCandidate struct {
	EntityID string       `json:"entityID"`
	Links    []ReviewLink `json:"links"`
}

// ReviewLink is a single pair of records that matched.
type ReviewLink struct {
	A     string       `json:"a"` // the submitted record, with version
	B     string       `json:"b"` // the already known record, with version
	Rules []ReviewRule `json:"rules"`
}

// ReviewRule is a rule that matched for a link, together with its score.
type ReviewRule struct {
	ID    string `json:"id"`
	Score uint8  `json:"score"`
}

// ReviewActorKind distinguishes the kinds of actors that can review a case.
type ReviewActorKind string

const (
	// ReviewActorHuman is a person reviewing through a user interface.
	ReviewActorHuman ReviewActorKind = "HUMAN"
	// ReviewActorAgent is an AI agent reviewing on behalf of a person.
	ReviewActorAgent ReviewActorKind = "AGENT"
	// ReviewActorProcess is an automated process applying earlier decisions.
	ReviewActorProcess ReviewActorKind = "PROCESS"
)

// ReviewActor identifies who claimed or resolved a case.
type ReviewActor struct {
	Kind ReviewActorKind `json:"kind"`
	ID   string          `json:"id"`   // account id, or a stable installation id
	Name string          `json:"name"` // for display only
}

// CreateReviewCaseInput identifies the entity to put under review.
type CreateReviewCaseInput struct {
	EntityID string `json:"entityID"`
}

// CreateReviewCaseOutput provides the queued case.
//
// Created reports whether the case was actually added. A case that was already
// waiting for this entity is returned unchanged instead of queueing a second
// one for the same thing.
type CreateReviewCaseOutput struct {
	Case    *ReviewCase `json:"case"`
	Created bool        `json:"created"`
}

// ClaimReviewCaseInput identifies the case to claim and who claims it.
type ClaimReviewCaseInput struct {
	ID    string      `json:"id"`
	Actor ReviewActor `json:"actor"`
}

// ClaimReviewCaseOutput provides the claimed case together with the entity that
// was locked for the review.
//
// EntityID is re-derived from the records of the case and may therefore differ
// from the entity ID stored on the case.
type ClaimReviewCaseOutput struct {
	Case     *ReviewCase `json:"case"`
	EntityID string      `json:"entityID"`
	Lock     string      `json:"lock"`
	Stale    bool        `json:"stale"` // the case no longer applies, nothing was locked
}

// ReleaseReviewCaseInput identifies the case to release.
//
// EntityID and Lock are optional. They are only required when the case itself no
// longer exists, which is the case when recovering a lock after a crash.
type ReleaseReviewCaseInput struct {
	ID       string `json:"id"`
	EntityID string `json:"entityID"`
	Lock     string `json:"lock"`
}

// ReleaseReviewCaseOutput reports whether a lock was actually released.
type ReleaseReviewCaseOutput struct {
	Released bool `json:"released"`
}

// ReviewLinkVerdict is the decision for a single link of a case.
//
// The link is addressed by the pair of records it connects, which is its
// natural key, and the order of the two does not matter.
type ReviewLinkVerdict struct {
	A    string `json:"a"`
	B    string `json:"b"`
	Keep bool   `json:"keep"`
}

// ResolveReviewCaseInput provides the verdicts for a claimed case.
type ResolveReviewCaseInput struct {
	ID       string              `json:"id"`
	Verdicts []ReviewLinkVerdict `json:"verdicts"`
	Actor    ReviewActor         `json:"actor"`
	Reason   string              `json:"reason"`
}

// ResolveReviewCaseOutput reports whether the resolution triggered a
// disassembly.
//
// The disassembly is asynchronous, a triggered disassembly has been enqueued but
// not necessarily applied yet.
type ResolveReviewCaseOutput struct {
	Triggered bool `json:"triggered"`
}

// ReviewDecisionsInput selects the decisions to return.
//
// The decisions are looked up by the records that were under review, not by the
// entity they belonged to: an entity ID may have merged away since the decision
// was made, a record ID is permanent.
type ReviewDecisionsInput struct {
	RecordIDs []string `json:"recordIDs"` // with or without version
	Cursor    *string  `json:"cursor"`    // opaque, empty starts at the newest
	Limit     *int     `json:"limit"`
}

// ReviewDecisionsOutput is the result of a ReviewDecisions call, newest
// decision first.
type ReviewDecisionsOutput struct {
	Decisions  []*ReviewDecision `json:"decisions"`
	NextCursor *string           `json:"nextCursor"` // nil means end of the history
}

// ReviewDecision is the permanent record of a resolved review case.
type ReviewDecision struct {
	ID        string              `json:"id"` // the ID of the case that was decided
	Case      *ReviewCase         `json:"case"`
	Verdicts  []ReviewLinkVerdict `json:"verdicts"`
	Actor     ReviewActor         `json:"actor"`
	Reason    string              `json:"reason"`
	DecidedAt time.Time           `json:"decidedAt"`
}
