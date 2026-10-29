package dispatcher

import (
	"encoding/json"
	"fmt"

	api "github.com/tilotech/tilores-plugin-api"
)

// AssembleEvent is used to parse a Kinesis Event that contains the data for
// the assemble lambda.
//
// The event can contain different types of payload depending on the action that
// needs to be performed.
//
// In case of type "ASSEMBLE", the payload contains []*api.Record entries.
//
// In case of type "DISASSEMBLE", the payload is a *dispatcher.DisassembleInput.
//
// In case of type "REVIEW_ASSEMBLE", the payload is a
// *dispatcher.ReviewAssembleInput.
//
// In case of type "REVIEW_CASE_RAISED", the payload is a
// *dispatcher.ReviewCaseRaisedInput.
//
// For backwards compatibility, the old event input for assemble requests is
// also supported. The output after unmarshaling will be the same as for type
// "ASSEMBLE".
type AssembleEvent struct {
	Type    string              `json:"type"`
	Payload any                 `json:"payload"`
	Flags   []AssembleEventFlag `json:"flags"`
}

const (
	// EventTypeAssemble is used when the payload is for the assemble process.
	EventTypeAssemble = "ASSEMBLE"

	// EventTypeDisassemble is used when the payload is for the disassemble process.
	EventTypeDisassemble = "DISASSEMBLE"

	// EventTypeReviewAssemble is used when the payload assembles the records of
	// a resolved review case again, see ReviewAssembleInput.
	EventTypeReviewAssemble = "REVIEW_ASSEMBLE"

	// EventTypeReviewCaseRaised is used after an assembly raised a review case,
	// see ReviewCaseRaisedInput.
	EventTypeReviewCaseRaised = "REVIEW_CASE_RAISED"
)

// ReviewCaseRaisedInput is the payload of a REVIEW_CASE_RAISED event. The
// assembly that raised a review case writes it once its plan released its
// locks. It is not meant to be sent by anybody else.
//
// If past decisions answer every link of the case - the newest decision with a
// verdict for the same pair of records, whatever rules matched it - the case
// is resolved with their verdicts, by an actor of kind PROCESS. Otherwise the
// event changes nothing and the case waits for a person.
type ReviewCaseRaisedInput struct {
	CaseID string `json:"caseID"`
}

// AssembleEventFlag defines additional flags that can be send in combination
// with an assemble event.
type AssembleEventFlag string

const (
	// StateRectification is an AssembleEventFlag used to indicate that the
	// provided event represents a correction of previously invalid, inconsistent,
	// or incomplete state within the entity graph.
	//
	// This flag signals that the event must trigger reprocessing of one or more
	// records to rectify a known incorrect state, typically caused by race
	// conditions or parallel processing during entity resolution.
	//
	// In rare cases, the invalid state may have already been corrected
	// automatically between the detection of the issue and the processing of the
	// flagged event. In such situations, the event will be ignored, since the
	// record involved already reflects the corrected state.
	//
	// It is only relevant if the event type is ASSEMBLE, and should be used
	// exclusively for automated corrections that restore canonical relationships
	// between entities.
	StateRectification AssembleEventFlag = "STATE_RECTIFICATION"
)

// UnmarshalJSON parses the provided bytes and populates the AssembleEvent.
func (r *AssembleEvent) UnmarshalJSON(b []byte) error {
	partial := &struct {
		Type    string
		Payload json.RawMessage
		Flags   []AssembleEventFlag
	}{}
	err := json.Unmarshal(b, partial)
	if err != nil {
		partial.Type = EventTypeAssemble
		partial.Payload = b
	}
	var payload any
	switch partial.Type {
	case EventTypeAssemble:
		pl := []*api.Record{}
		err = json.Unmarshal(partial.Payload, &pl)
		payload = pl
	case EventTypeDisassemble:
		payload = &DisassembleInput{}
		err = json.Unmarshal(partial.Payload, payload)
	case EventTypeReviewAssemble:
		payload = &ReviewAssembleInput{}
		err = json.Unmarshal(partial.Payload, payload)
	case EventTypeReviewCaseRaised:
		payload = &ReviewCaseRaisedInput{}
		err = json.Unmarshal(partial.Payload, payload)
	default:
		return fmt.Errorf("invalid type %s", partial.Type)
	}
	if err != nil {
		return err
	}
	r.Type = partial.Type
	r.Payload = payload
	r.Flags = partial.Flags
	return nil
}
