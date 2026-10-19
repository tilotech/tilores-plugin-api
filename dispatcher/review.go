package dispatcher

import (
	"time"

	api "github.com/tilotech/tilores-plugin-api"
)

// ReviewCaseStatus describes where a review case currently is in its lifecycle.
type ReviewCaseStatus string

const (
	// ReviewCaseOpen marks a case that nobody is working on.
	ReviewCaseOpen ReviewCaseStatus = "OPEN"
	// ReviewCaseClaimed marks a case that an actor is currently reviewing and
	// for which the affected entities are locked.
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

// ReviewCase describes an assembly that found one or more entities on evidence
// that was too weak to be trusted.
//
// Depending on the join strategy of the configuration, the assembly either
// joined those entities and the review decides whether the join stands, or it
// deferred the join and the review decides whether it happens. Which one
// applies is a property of every single candidate, see ReviewCandidate.
type ReviewCase struct {
	ID         string            `json:"id"`
	Status     ReviewCaseStatus  `json:"status"`
	CreatedAt  time.Time         `json:"createdAt"`
	RecordIDs  []string          `json:"recordIDs"` // the submitted records, with version
	EntityID   string            `json:"entityID"`  // the entity of the submitted records at detection time
	Candidates []ReviewCandidate `json:"candidates"`
	ClaimedBy  *ReviewActor      `json:"claimedBy"`
	ClaimedAt  *time.Time        `json:"claimedAt"`
}

// ReviewCandidate is a single entity that the submitted records matched too
// weakly.
//
// Deferred reports that the join with the candidate was deferred: the
// candidate was kept apart from the submitted records instead of being joined.
// A case with deferred candidates may still list joined ones, e.g. the entity
// the submitted records already belonged to, whose join is never deferred.
type ReviewCandidate struct {
	EntityID string       `json:"entityID"`
	Links    []ReviewLink `json:"links"`
	Deferred bool         `json:"deferred"`
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
//
// This is the manual counterpart to detection and is never called from an
// assembly: it exists for the case a person opens after spotting a wrong merge
// that the criterion never flagged. The entity ID is therefore all there is to
// give - at that point nobody knows yet which of the entity's links are the
// mistake, which is exactly what the review is going to establish - so the
// resulting case carries no links and a verdict may name any pair of its
// records.
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

// ClaimReviewCaseOutput provides the claimed case together with the entities
// that were locked for the review.
//
// EntityID is the entity of the submitted records. It is re-derived from the
// records of the case and may therefore differ from the entity ID stored on the
// case.
//
// CandidateEntityIDs are the entities of the deferred candidates, also
// re-derived from their records. They are locked with the same lock, so that
// nothing the review is about can change while it is decided. The list is
// empty for a case without deferred candidates.
type ClaimReviewCaseOutput struct {
	Case               *ReviewCase `json:"case"`
	EntityID           string      `json:"entityID"`
	CandidateEntityIDs []string    `json:"candidateEntityIDs"`
	Lock               string      `json:"lock"`
	Stale              bool        `json:"stale"` // the case no longer applies, nothing was locked
}

// ReleaseReviewCaseInput identifies the case to release.
//
// EntityID, CandidateEntityIDs and Lock are optional. They are only required
// when the case itself no longer exists, which is the case when recovering a
// lock after a crash.
type ReleaseReviewCaseInput struct {
	ID                 string   `json:"id"`
	EntityID           string   `json:"entityID"`
	CandidateEntityIDs []string `json:"candidateEntityIDs"`
	Lock               string   `json:"lock"`
}

// ReleaseReviewCaseOutput reports whether a lock was actually released.
type ReleaseReviewCaseOutput struct {
	Released bool `json:"released"`
}

// ReviewLinkVerdict is the decision for a single link of a case.
//
// The link is addressed by the pair of records it connects, which is its
// natural key, and the order of the two does not matter.
//
// Keep tells whether the two records legitimately belong to the same entity.
// For a link that was joined, keeping it leaves it in place and rejecting it
// removes it. For a link of a deferred candidate, keeping it joins the two
// records and rejecting it leaves them apart.
type ReviewLinkVerdict struct {
	A    string `json:"a"`
	B    string `json:"b"`
	Keep bool   `json:"keep"`
}

// ResolveReviewCaseInput provides the verdicts for a claimed case.
//
// Verdicts must cover every link the case carries. A link without a verdict is
// refused rather than defaulted to being kept, because silently confirming a
// join that nobody looked at is the one outcome this feature exists to prevent.
// Further verdicts are allowed and are how a review concludes that the flagged
// link is fine and a neighbouring one is the actual mistake: a link is named by
// the pair of records it connects, so such a pair needs nothing the case would
// have to know about beforehand. For a case with deferred candidates, such a
// pair must connect one submitted record with one record that was not
// submitted, because the records are applied by assembling the submitted ones
// again, which can only decide about the links of those.
//
// Lock is the lock ClaimReviewCase handed out. It is a generation marker on the
// claim rather than a capability: a claim that was released and taken again
// mints a new value, so a caller resolving after losing its claim is refused
// instead of deciding on the current holder's behalf, which the actor alone
// cannot catch when the new holder is the same actor. It is only ever compared
// against the claim - the disassembly or assembly that applies the verdicts
// takes over the lock the case itself recorded.
//
// It is not a secret: a listed case carries the claiming actor and the moment
// of the claim, which is what the value is made of. Nor does it separate two
// concurrent sessions of one actor. The lock model is cooperative throughout -
// ReleaseReviewCase takes no lock at all, so that a claim whose holder is gone
// can always be freed.
type ResolveReviewCaseInput struct {
	ID       string              `json:"id"`
	Verdicts []ReviewLinkVerdict `json:"verdicts"`
	Actor    ReviewActor         `json:"actor"`
	Reason   string              `json:"reason"`
	Lock     string              `json:"lock"`
}

// ResolveReviewCaseOutput reports whether the resolution triggered a change.
//
// For a case that joined its candidates, the change is the disassembly of the
// rejected links. For a case with deferred candidates, it is the assembly that
// joins the kept links. Either is asynchronous: a triggered change has been
// enqueued but not necessarily applied yet. The locks of the case are released
// once it has been applied, or right away if nothing was triggered.
type ResolveReviewCaseOutput struct {
	Triggered bool `json:"triggered"`
}

// ReviewAssembleInput is the payload of a REVIEW_ASSEMBLE event, which applies
// the verdicts of a review case with deferred candidates.
//
// The records of the case are assembled again, with every kept link counting
// as approved and every rejected link being ignored, so that the current rules
// and consistency decide about the result. It is written by ResolveReviewCase
// and not meant to be sent by anybody else.
//
// Records carry the version they are stored with. Their data is unchanged, a
// new version is created nevertheless. The assembly is dropped if a newer
// version exists or a record no longer exists by then.
//
// The record IDs of the verdicts always carry a version.
//
// Lock is the review lock of the claim. The assembly takes it over from every
// entity that still carries it, including those of LockedEntityIDs that it does
// not otherwise touch, and releases all of them once it has been applied.
type ReviewAssembleInput struct {
	CaseID          string              `json:"caseID"`
	Records         []*api.Record       `json:"records"`
	Verdicts        []ReviewLinkVerdict `json:"verdicts"`
	Lock            string              `json:"lock"`
	LockedEntityIDs []string            `json:"lockedEntityIDs"`
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
