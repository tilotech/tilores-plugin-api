package dispatcher

import (
	"context"

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
	Submit(ctx context.Context, records []*api.Record) (*SubmissionResult, error)
	Search(ctx context.Context, parameters *api.SearchParameters) ([]*api.Entity, error)
	RemoveConnectionBan(ctx context.Context, input *RemoveConnectionBanInput) error
}

// SubmissionResult provides additional information about a successful
// data submission.
type SubmissionResult struct {
	RecordsAdded int
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
