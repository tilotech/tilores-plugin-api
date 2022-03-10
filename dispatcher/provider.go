package dispatcher

import (
	"context"
	"fmt"

	"github.com/tilotech/go-plugin"
)

// Provide returns the plugin.Provider for the given Dispatcher.
func Provide(impl Dispatcher) plugin.Provider {
	return &provider{
		impl: impl,
	}
}

type provider struct {
	impl Dispatcher
}

const (
	entityMethod              = "/entity"
	submitMethod              = "/submit"
	disassembleMethod         = "/disassemble"
	removeConnectionBanMethod = "/removeconnectionban"
	searchMethod              = "/search"
)

func (p *provider) Provide(method string) (plugin.RequestParameter, plugin.InvokeFunc, error) {
	switch method {
	case entityMethod:
		return &EntityInput{}, p.Entity, nil
	case submitMethod:
		return &SubmitInput{}, p.Submit, nil
	case disassembleMethod:
		return &DisassembleInput{}, p.Disassemble, nil
	case removeConnectionBanMethod:
		return &RemoveConnectionBanInput{}, p.RemoveConnectionBan, nil
	case searchMethod:
		return &SearchInput{}, p.Search, nil
	}
	return nil, nil, fmt.Errorf("invalid method %v", method)
}

func (p *provider) Entity(ctx context.Context, params plugin.RequestParameter) (interface{}, error) {
	return p.impl.Entity(ctx, params.(*EntityInput))
}

func (p *provider) Submit(ctx context.Context, params plugin.RequestParameter) (interface{}, error) {
	return p.impl.Submit(ctx, params.(*SubmitInput))
}

func (p *provider) Disassemble(ctx context.Context, params plugin.RequestParameter) (interface{}, error) {
	return p.impl.Disassemble(ctx, params.(*DisassembleInput))
}

func (p *provider) RemoveConnectionBan(ctx context.Context, params plugin.RequestParameter) (interface{}, error) {
	return nil, p.impl.RemoveConnectionBan(ctx, params.(*RemoveConnectionBanInput))
}

func (p *provider) Search(ctx context.Context, params plugin.RequestParameter) (interface{}, error) {
	return p.impl.Search(ctx, params.(*SearchInput))
}
