package dispatcher

import (
	"context"
	"net/rpc"

	api "github.com/tilotech/tilores-plugin-api"
)

type client struct {
	client *rpc.Client
}

func (c *client) Submit(ctx context.Context, input *SubmitInput) (*SubmitOutput, error) {
	var submitOutput SubmitOutput
	err := c.client.Call(
		"Plugin.Submit",
		map[string]interface{}{
			"input": input,
		},
		&submitOutput,
	)
	if err != nil {
		return nil, err
	}
	return &submitOutput, nil
}

func (c *client) Disassemble(ctx context.Context, input *DisassembleInput) (*DisassembleOutput, error) {
	var output DisassembleOutput
	err := c.client.Call(
		"Plugin.Disassemble",
		map[string]interface{}{
			"input": input,
		},
		&output,
	)
	return &output, err
}

func (c *client) RemoveConnectionBan(ctx context.Context, input *RemoveConnectionBanInput) error {
	var reply interface{}
	err := c.client.Call(
		"Plugin.RemoveConnectionBan",
		map[string]interface{}{
			"input": input,
		},
		&reply,
	)
	return err
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

func (c *client) Search(ctx context.Context, parameters *api.SearchParameters) ([]*api.Entity, error) {
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
