package dispatcher

import (
	"context"
	"net/rpc"
)

type client struct {
	client *rpc.Client
}

func (c *client) Submit(ctx context.Context, input *SubmitInput) (*SubmitOutput, error) {
	args := map[string]interface{}{
		"input": input,
	}
	if deadline, ok := ctx.Deadline(); ok {
		args["deadline"] = deadline
	}
	var submitOutput SubmitOutput
	err := c.client.Call(
		"Plugin.Submit",
		args,
		&submitOutput,
	)
	if err != nil {
		return nil, err
	}
	return &submitOutput, nil
}

func (c *client) Disassemble(ctx context.Context, input *DisassembleInput) (*DisassembleOutput, error) {
	args := map[string]interface{}{
		"input": input,
	}
	if deadline, ok := ctx.Deadline(); ok {
		args["deadline"] = deadline
	}
	var output DisassembleOutput
	err := c.client.Call(
		"Plugin.Disassemble",
		args,
		&output,
	)
	return &output, err
}

func (c *client) RemoveConnectionBan(ctx context.Context, input *RemoveConnectionBanInput) error {
	args := map[string]interface{}{
		"input": input,
	}
	if deadline, ok := ctx.Deadline(); ok {
		args["deadline"] = deadline
	}
	var reply interface{}
	err := c.client.Call(
		"Plugin.RemoveConnectionBan",
		args,
		&reply,
	)
	return err
}

func (c *client) Entity(ctx context.Context, input *EntityInput) (*EntityOutput, error) {
	args := map[string]interface{}{
		"input": input,
	}
	if deadline, ok := ctx.Deadline(); ok {
		args["deadline"] = deadline
	}
	var output EntityOutput
	err := c.client.Call(
		"Plugin.Entity",
		args,
		&output,
	)
	if err != nil {
		return nil, err
	}
	return &output, nil
}

func (c *client) Search(ctx context.Context, input *SearchInput) (*SearchOutput, error) {
	args := map[string]interface{}{
		"input": input,
	}
	if deadline, ok := ctx.Deadline(); ok {
		args["deadline"] = deadline
	}
	var output SearchOutput
	err := c.client.Call(
		"Plugin.Search",
		args,
		&output,
	)
	if err != nil {
		return nil, err
	}
	return &output, nil
}
