package dispatcher

import "time"

// EntityStreamEventType describes what happened to the entities of an event.
type EntityStreamEventType string

const (
	// EntityStreamEventCreate means that a new entity was created.
	EntityStreamEventCreate EntityStreamEventType = "CREATE"
	// EntityStreamEventUpdate means that an entity has been updated.
	EntityStreamEventUpdate EntityStreamEventType = "UPDATE"
	// EntityStreamEventMerge means that two entities have been merged.
	EntityStreamEventMerge EntityStreamEventType = "MERGE"
	// EntityStreamEventSplit means that one entity has been split into two or more
	// entities.
	EntityStreamEventSplit EntityStreamEventType = "SPLIT"
	// EntityStreamEventDelete means that an entity was deleted.
	EntityStreamEventDelete EntityStreamEventType = "DELETE"
	// EntityStreamEventOffload means that the event was too large for the stream.
	// It has no data, but a URL for downloading the original event.
	EntityStreamEventOffload EntityStreamEventType = "OFFLOAD"
)

// EntityStreamEventsInput defines how many events to read from the entity
// stream.
type EntityStreamEventsInput struct {
	Limit int `json:"limit"`
}

// EntityStreamEventsOutput is the result of an EntityStreamEvents call.
//
// Fewer events than requested do not mean that the stream is empty.
type EntityStreamEventsOutput struct {
	Events []*EntityStreamEvent `json:"events"`
}

// EntityStreamEvent is a single event from the entity stream.
//
// Apart from the AckToken, it has the same format as the events of the stream
// itself.
type EntityStreamEvent struct {
	// AckToken is opaque and identifies this delivery of the event. An event that
	// is delivered again has a different token.
	AckToken  string                 `json:"ackToken"`
	Type      EntityStreamEventType  `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Data      *EntityStreamEventData `json:"data,omitempty"`
	URL       *string                `json:"url,omitempty"`
}

// EntityStreamEventData contains the affected entities of an EntityStreamEvent.
type EntityStreamEventData struct {
	Entities         []EntityStreamEntityInfo `json:"entities"`
	ObsoleteEntities []EntityStreamEntityInfo `json:"obsoleteEntities"`
}

// EntityStreamEntityInfo contains the information about one affected entity.
type EntityStreamEntityInfo struct {
	ID                string         `json:"id"`
	RecordIDs         []string       `json:"recordIDs"`
	ObsoleteRecordIDs []string       `json:"obsoleteRecordIDs"`
	Insights          map[string]any `json:"insights"`
}

// AcknowledgeEntityStreamEventsInput lists the events to remove from the entity
// stream, using the tokens from an EntityStreamEvents call.
type AcknowledgeEntityStreamEventsInput struct {
	AckTokens []string `json:"ackTokens"`
}

// AcknowledgeEntityStreamEventsOutput is the result of an
// AcknowledgeEntityStreamEvents call.
type AcknowledgeEntityStreamEventsOutput struct {
	// Acknowledged is the number of accepted tokens. An accepted token of an earlier
	// delivery does not always remove the event, it is then delivered again.
	Acknowledged int `json:"acknowledged"`

	// InvalidAckTokens lists the tokens that were rejected. Repeating the call does
	// not help for them. An event that is still in the stream is delivered again
	// with a new token.
	InvalidAckTokens []string `json:"invalidAckTokens"`
}
