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
	Entity(ctx context.Context, id string) (*api.Entity, error)
	Submit(ctx context.Context, input *SubmitInput) (*SubmitOutput, error)
	Search(ctx context.Context, parameters *api.SearchParameters) ([]*api.Entity, error)
	Disassemble(ctx context.Context, input *DisassembleInput) (*DisassembleOutput, error)
	RemoveConnectionBan(ctx context.Context, input *RemoveConnectionBanInput) error
}

// SubmitInput includes the data requied to submit
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
	gob.Register(&SubmitInput{})
	gob.Register(&DisassembleInput{})
	gob.Register(&RemoveConnectionBanInput{})
}
