package dispatcher

import (
	"context"
	"net/rpc"

	api "gitlab.com/tilotech/tilores-plugin-api"
)

type client struct {
	client *rpc.Client
}

func (c *client) Entity(ctx context.Context, id string) (*api.Entity, error) {
	var entity api.Entity
	err := c.client.Call(
		"Plugin.Entity",
		map[string]interface{}{
			"id": id,
		},
		&entity,
	)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (c *client) Submit(ctx context.Context, records []*api.Record) (*SubmissionResult, error) {
	var submissionResult SubmissionResult
	err := c.client.Call(
		"Plugin.Submit",
		map[string]interface{}{
			"records": records,
		},
		&submissionResult,
	)
	if err != nil {
		return nil, err
	}
	return &submissionResult, nil
}

func (c *client) Search(ctx context.Context, parameters map[string]interface{}) ([]*api.Entity, error) {
	var searchResult []*api.Entity
	err := c.client.Call(
		"Plugin.Search",
		map[string]interface{}{
			"parameters": parameters,
		},
		&searchResult,
	)
	if err != nil {
		return nil, err
	}
	return searchResult, nil
}
