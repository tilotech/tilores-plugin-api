package dispatcher

import (
	"context"
	"fmt"
	"os"

	"github.com/tilotech/go-plugin"
)

// Connect starts the dispatcher plugin and returns a proxy that implements the
// dispatcher interface.
func Connect(starter plugin.Starter, config *plugin.Config) (Dispatcher, plugin.TermFunc, error) {
	client, term, err := plugin.Start(starter, fmt.Sprintf("%v/dispatcher", os.TempDir()), config)
	if err != nil {
		return nil, nil, err
	}
	return &proxy{
		client: client,
	}, term, nil
}

type proxy struct {
	client *plugin.Client
}

func (p *proxy) Entity(ctx context.Context, input *EntityInput) (*EntityOutput, error) {
	response := &EntityOutput{}
	err := p.client.Call(ctx, entityMethod, input, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (p *proxy) Submit(ctx context.Context, input *SubmitInput) (*SubmitOutput, error) {
	response := &SubmitOutput{}
	err := p.client.Call(ctx, submitMethod, input, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (p *proxy) Disassemble(ctx context.Context, input *DisassembleInput) (*DisassembleOutput, error) {
	response := &DisassembleOutput{}
	err := p.client.Call(ctx, disassembleMethod, input, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (p *proxy) RemoveConnectionBan(ctx context.Context, input *RemoveConnectionBanInput) error {
	var response interface{}
	return p.client.Call(ctx, removeConnectionBanMethod, input, response)
}

func (p *proxy) Search(ctx context.Context, input *SearchInput) (*SearchOutput, error) {
	response := &SearchOutput{}
	err := p.client.Call(ctx, searchMethod, input, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
