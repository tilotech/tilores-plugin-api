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
	Search(ctx context.Context, parameters map[string]interface{}) ([]*api.Entity, error)
}

// SubmissionResult provides additional information about a successful
// data submission.
type SubmissionResult struct {
	RecordsAdded int
}
