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
//
// RecordIDs are only set for the other entity of a case that was created by
// hand for two entities, see CreateReviewCaseInput. They are the records of
// that entity, with version, and are how its current entity is found again.
// Such a candidate has no links: no rule connects the two entities, and keeping
// a link between one of its records and a record of the case joins them.
type ReviewCandidate struct {
	EntityID  string       `json:"entityID"`
	Links     []ReviewLink `json:"links"`
	Deferred  bool         `json:"deferred"`
	RecordIDs []string     `json:"recordIDs"`
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
// records. The case is about the records the entity holds when it is created:
// if the entity splits before the case is claimed, the claim locks all its
// parts, and keeping a link between two of them joins them again.
//
// OtherEntityID is optional and names a second entity that no rule connects
// with the first one, for the case a person opens after spotting that the two
// are the same. The case then lists the second entity as a candidate with its
// records, and keeping a link between a record of each joins the two entities.
// The order of the two entities does not matter.
type CreateReviewCaseInput struct {
	EntityID      string `json:"entityID"`
	OtherEntityID string `json:"otherEntityID"`
}

// CreateReviewCaseOutput provides the queued case.
//
// Created reports whether the case was actually added. A case that was already
// waiting for these entities is returned unchanged instead of queueing a second
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
// The records of a case are its submitted records, both records of each of its
// links, and the records of the other entity of a case created for two
// entities. The claim locks every entity they are in now, whatever kind the
// case is, so that a review can change anything between them and nothing it is
// about can change while it is decided. An entity that split since the case was
// created is locked in all its parts.
//
// EntityID is the entity of the first submitted record that still exists. It
// is re-derived from the records and may therefore differ from the entity ID
// stored on the case. CandidateEntityIDs are all further entities, locked with
// the same lock.
//
// A case is stale once its submitted records are all gone. A case with deferred
// candidates is also stale once a submitted record has a newer version or was
// deleted, and a case whose candidates are all deferred once no deferred link
// connects two existing records in different entities any more: other links
// joined them, or the records of the candidate are gone. A case created for two
// entities is stale once the records of the other entity are all gone or all
// its records share one entity. A stale case is removed from the queue.
//
// A claim fails with the error type tilores.entity_busy while one of these
// entities is locked by another operation: an ingest, the assembly of a decision
// that was just applied, or the claim of another reviewer. Nothing is locked
// then, and the claim may succeed a moment later.
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
// The same rules hold for every kind of case: rejecting a pair that an edge
// connects removes every edge between the two, and keeping a pair that no edge
// connects adds one, which joins their entities if they are two. Rejecting an
// unconnected pair and keeping a connected one change nothing.
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
// Further verdicts are allowed and are how a review removes or creates an edge
// that the case does not name: a link is named by the pair of records it
// connects, so such a pair needs nothing the case would have to know about
// beforehand. A further pair must connect two different records that exist and
// belong to the entities the claim locked. A link of the case itself may name a
// record that was deleted since; its verdict then changes nothing.
//
// The verdicts are applied the same way for every kind of case, see
// ReviewLinkVerdict and ReviewAssembleInput. A kept edge that joins two
// entities is not applied if the consistency rules keep them apart; edges that
// join separate groups of entities are checked group by group.
//
// Lock is the lock ClaimReviewCase handed out. It is a generation marker on the
// claim rather than a capability: a claim that was released and taken again
// mints a new value, so a caller resolving after losing its claim is refused
// instead of deciding on the current holder's behalf, which the actor alone
// cannot catch when the new holder is the same actor. It is only ever compared
// against the claim - the assembly that applies the verdicts takes over the
// lock the case itself recorded.
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
// The change is the assembly that applies the verdicts, see
// ReviewAssembleInput. It is asynchronous: a triggered change has been enqueued
// but not necessarily applied yet. The locks of the case are released once it
// has been applied, or right away if no verdict changes anything.
type ResolveReviewCaseOutput struct {
	Triggered bool `json:"triggered"`
}

// ReviewAssembleInput is the payload of a REVIEW_ASSEMBLE event, which applies
// the verdicts of a review case. It is written by ResolveReviewCase and not
// meant to be sent by anybody else.
//
// The verdicts are applied to the graphs of the entities the claim locked, in
// one assembly. A rejected pair that an edge connects loses every edge between
// the two records, which splits what is no longer connected. A kept pair that
// no edge connects gets an edge: with the rules of the link if Links names the
// pair, and as a manual link otherwise. A kept pair across two entities joins
// them, unless the consistency rules keep them apart. No record is assembled
// again and no record version changes.
//
// The record IDs of the verdicts always carry a version. Links are the links of
// the case, with the rules that matched them.
//
// Lock is the review lock of the claim. The assembly takes it over from every
// entity that still carries it, including those of LockedEntityIDs that it does
// not otherwise touch, and releases all of them once it has been applied. Only
// the entities that still carry the review lock take part, so that an event
// that is delivered again after the locks were released changes nothing.
//
// Records is empty in every event written now. An event that an earlier
// version wrote for a case with deferred candidates carries its submitted
// records instead: they are assembled again with every kept link counting as
// approved, and a new version is created. Such an event is dropped if a newer
// version exists or a record no longer exists by then.
type ReviewAssembleInput struct {
	CaseID          string              `json:"caseID"`
	Records         []*api.Record       `json:"records"`
	Verdicts        []ReviewLinkVerdict `json:"verdicts"`
	Links           []ReviewLink        `json:"links"`
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
