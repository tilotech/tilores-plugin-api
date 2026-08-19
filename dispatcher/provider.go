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
	entityByRecordMethod      = "/entity-by-record"
	submitMethod              = "/submit"
	submitWithPreviewMethod   = "/submit-with-preview"
	disassembleMethod         = "/disassemble"
	removeConnectionBanMethod = "/removeconnectionban"
	searchMethod              = "/search"
	assemblyStatusMethod      = "/assembly-status"
	reviewCasesMethod         = "/review-cases"
	claimReviewCaseMethod     = "/claim-review-case"
	releaseReviewCaseMethod   = "/release-review-case"
	resolveReviewCaseMethod   = "/resolve-review-case"
)

func (p *provider) Provide(method string) (plugin.RequestParameter, plugin.InvokeFunc, error) {
	switch method {
	case entityMethod:
		return &EntityInput{}, p.Entity, nil
	case entityByRecordMethod:
		return &EntityByRecordInput{}, p.EntityByRecord, nil
	case submitMethod:
		return &SubmitInput{}, p.Submit, nil
	case submitWithPreviewMethod:
		return &SubmitWithPreviewInput{}, p.SubmitWithPreview, nil
	case disassembleMethod:
		return &DisassembleInput{}, p.Disassemble, nil
	case removeConnectionBanMethod:
		return &RemoveConnectionBanInput{}, p.RemoveConnectionBan, nil
	case searchMethod:
		return &SearchInput{}, p.Search, nil
	case assemblyStatusMethod:
		return nil, p.AssemblyStatus, nil
	case reviewCasesMethod:
		return &ReviewCasesInput{}, p.ReviewCases, nil
	case claimReviewCaseMethod:
		return &ClaimReviewCaseInput{}, p.ClaimReviewCase, nil
	case releaseReviewCaseMethod:
		return &ReleaseReviewCaseInput{}, p.ReleaseReviewCase, nil
	case resolveReviewCaseMethod:
		return &ResolveReviewCaseInput{}, p.ResolveReviewCase, nil
	}
	return nil, nil, fmt.Errorf("invalid method %v", method)
}

func (p *provider) Entity(ctx context.Context, params plugin.RequestParameter) (any, error) {
	return p.impl.Entity(ctx, params.(*EntityInput))
}

func (p *provider) EntityByRecord(ctx context.Context, params plugin.RequestParameter) (any, error) {
	return p.impl.EntityByRecord(ctx, params.(*EntityByRecordInput))
}

func (p *provider) Submit(ctx context.Context, params plugin.RequestParameter) (any, error) {
	return p.impl.Submit(ctx, params.(*SubmitInput))
}

func (p *provider) SubmitWithPreview(ctx context.Context, params plugin.RequestParameter) (any, error) {
	return p.impl.SubmitWithPreview(ctx, params.(*SubmitWithPreviewInput))
}

func (p *provider) Disassemble(ctx context.Context, params plugin.RequestParameter) (any, error) {
	return p.impl.Disassemble(ctx, params.(*DisassembleInput))
}

func (p *provider) RemoveConnectionBan(ctx context.Context, params plugin.RequestParameter) (any, error) {
	return nil, p.impl.RemoveConnectionBan(ctx, params.(*RemoveConnectionBanInput))
}

func (p *provider) Search(ctx context.Context, params plugin.RequestParameter) (any, error) {
	return p.impl.Search(ctx, params.(*SearchInput))
}

func (p *provider) AssemblyStatus(ctx context.Context, _ plugin.RequestParameter) (any, error) {
	return p.impl.AssemblyStatus(ctx)
}

func (p *provider) ReviewCases(ctx context.Context, params plugin.RequestParameter) (any, error) {
	return p.impl.ReviewCases(ctx, params.(*ReviewCasesInput))
}

func (p *provider) ClaimReviewCase(ctx context.Context, params plugin.RequestParameter) (any, error) {
	return p.impl.ClaimReviewCase(ctx, params.(*ClaimReviewCaseInput))
}

func (p *provider) ReleaseReviewCase(ctx context.Context, params plugin.RequestParameter) (any, error) {
	return p.impl.ReleaseReviewCase(ctx, params.(*ReleaseReviewCaseInput))
}

func (p *provider) ResolveReviewCase(ctx context.Context, params plugin.RequestParameter) (any, error) {
	return p.impl.ResolveReviewCase(ctx, params.(*ResolveReviewCaseInput))
}
