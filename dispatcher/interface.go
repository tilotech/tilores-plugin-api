package dispatcher

import (
	"context"
	"encoding/gob"
	"time"

	api "github.com/tilotech/tilores-plugin-api"
)

// Dispatcher is the interface used for communicating between the public facing
// webserver API (typically GraphQL) and the internal TiloRes API.
//
// This interface is mostly created to support different deployment types, e.g.
// a local deployment with a fake implementation and a serverless deployment
// with the actual implementation.
//
// However, it might also offer the possibility to add data modifications on the
// customers side at a central place.
type Dispatcher interface {
	Entity(ctx context.Context, input *EntityInput) (*EntityOutput, error)
	Submit(ctx context.Context, input *SubmitInput) (*SubmitOutput, error)
	Search(ctx context.Context, input *SearchInput) (*SearchOutput, error)
	Disassemble(ctx context.Context, input *DisassembleInput) (*DisassembleOutput, error)
	RemoveConnectionBan(ctx context.Context, input *RemoveConnectionBanInput) error
}

// EntityInput includes the data required to get an entity by its ID
type EntityInput struct {
	ID string `json:"id"`
}

// EntityOutput the output of Entity call
type EntityOutput struct {
	Entity *api.Entity `json:"entity"`
}

// SearchInput includes the search parameters
type SearchInput struct {
	Parameters *api.SearchParameters `json:"parameters"`
}

// SearchOutput the output of Search call
type SearchOutput struct {
	Entities []*api.Entity `json:"entities"`
	Hits     Hits          `json:"hits"`
}

// Hits lists all matched rules per matched record id
//
// Example (in JSON):
// {
//   "550e8400-e29b-11d4-a716-446655440000": ["RULE-1", "RULE-2"],
//   "6ba7b810-9dad-11d1-80b4-00c04fd430c8": ["RULE-2"],
// }
type Hits map[string][]string

// IDs returns the record ids of the hits
func (h Hits) IDs() []string {
	ids := make([]string, len(h))
	i := 0
	for k := range h {
		ids[i] = k
		i++
	}
	return ids
}

// SubmitInput includes the data required to submit
type SubmitInput struct {
	Records []*api.Record `json:"records"`
}

// SubmitOutput provides additional information about a successful
// data submission.
type SubmitOutput struct {
	RecordsAdded int `json:"recordsAdded"`
}

// DisassembleInput is the data required to remove one or more edges or even records
//
// The metadata is required when disassemble is triggered by a real person,
// Otherwise it MAY be omitted.
type DisassembleInput struct {
	Reference           string            `json:"reference"`
	Edges               []DisassembleEdge `json:"edges"`
	RecordIDs           []string          `json:"recordIDs"`
	CreateConnectionBan bool              `json:"createConnectionBan"`
	Meta                DisassembleMeta   `json:"meta"`
	Timeout             *time.Duration    `json:"timeout"`
}

// DisassembleEdge represents a single edge to be removed
type DisassembleEdge struct {
	A string `json:"a"`
	B string `json:"b"`
}

// DisassembleMeta provides information who and why disassemble was started
type DisassembleMeta struct {
	User   string `json:"user"`
	Reason string `json:"reason"`
}

// DisassembleOutput informs about removed records and edges as well as the
// remaining entity ids
type DisassembleOutput struct {
	DeletedEdges   int32    `json:"deletedEdges"`
	DeletedRecords int32    `json:"deletedRecords"`
	EntityIDs      []string `json:"entityIDs"`
}

// RemoveConnectionBanInput contains the data required to remove a connection ban
type RemoveConnectionBanInput struct {
	Reference string                  `json:"reference"`
	EntityID  string                  `json:"entityID"`
	Others    []string                `json:"others"`
	Meta      RemoveConnectionBanMeta `json:"meta"`
}

// RemoveConnectionBanMeta provides information who and why the connection ban was removed
type RemoveConnectionBanMeta struct {
	User   string `json:"user"`
	Reason string `json:"reason"`
}

func init() {
	gob.Register(&EntityInput{})
	gob.Register(&SearchInput{})
	gob.Register(&SubmitInput{})
	gob.Register(&DisassembleInput{})
	gob.Register(&RemoveConnectionBanInput{})
}
