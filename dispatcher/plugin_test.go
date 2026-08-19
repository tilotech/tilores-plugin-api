package dispatcher_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tilotech/go-plugin"
	api "github.com/tilotech/tilores-plugin-api"
	"github.com/tilotech/tilores-plugin-api/dispatcher"
)

func TestPlugin(t *testing.T) {
	pluginImpl := &testDispatcher{}
	dsp, term, err := dispatcher.Connect(
		plugin.StartWithProvider(dispatcher.Provide(pluginImpl)),
		plugin.DefaultConfig(),
	)
	require.NoError(t, err)
	defer term()

	contextWithDeadline, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	entityOutput, err := dsp.Entity(contextWithDeadline, &dispatcher.EntityInput{ID: "abcd"})
	assert.NoError(t, err)
	require.NotNil(t, entityOutput)
	assert.Equal(t, 1, len(entityOutput.Entity.Records))
	assert.Equal(t, "bar", entityOutput.Entity.Records[0].Data["foo"])
	assert.Equal(t, 1, len(entityOutput.Entity.Edges))
	assert.Equal(t, 1, len(entityOutput.Entity.Duplicates))
	assert.True(t, pluginImpl.deadlineExists)

	entityByRecordOutput, err := dsp.EntityByRecord(contextWithDeadline, &dispatcher.EntityByRecordInput{ID: "12345"})
	assert.NoError(t, err)
	require.NotNil(t, entityByRecordOutput)
	assert.Equal(t, entityOutput.Entity, entityByRecordOutput.Entity)

	parameters := &api.SearchParameters{
		"foo": "bar",
	}
	searchOutput, err := dsp.Search(context.Background(), &dispatcher.SearchInput{Parameters: parameters})
	assert.NoError(t, err)
	assert.NotNil(t, searchOutput)
	require.Equal(t, 1, len(searchOutput.Entities))

	submitOutput, err := dsp.Submit(context.Background(), &dispatcher.SubmitInput{
		Records: []*api.Record{
			{
				ID: "12345",
				Data: map[string]any{
					"foo": "bar",
				},
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, submitOutput.RecordsAdded)

	submitWithPreviewOutput, err := dsp.SubmitWithPreview(context.Background(), &dispatcher.SubmitWithPreviewInput{
		Records: []*api.Record{
			{
				ID: "12345",
				Data: map[string]any{
					"foo": "bar",
				},
			},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, submitWithPreviewOutput)
	assert.Len(t, submitWithPreviewOutput.Entities, 1)

	disassembleOutput, err := dsp.Disassemble(context.Background(), &dispatcher.DisassembleInput{
		Edges: []dispatcher.DisassembleEdge{
			{
				A: "abc",
				B: "def",
			},
		},
		RecordIDs: []string{
			"12345",
		},
		CreateConnectionBan: true,
		Meta: &dispatcher.DisassembleMeta{
			User:   "someUser",
			Reason: "someReason",
		},
		Lock: "review:someUser:1",
	})
	assert.NoError(t, err)
	assert.True(t, disassembleOutput.Triggered)
	assert.Equal(t, "review:someUser:1", pluginImpl.disassembleLock)

	err = dsp.RemoveConnectionBan(context.Background(), &dispatcher.RemoveConnectionBanInput{
		Reference: "123123",
		EntityID:  "someID",
		Others:    []string{"someOtherID"},
		Meta: dispatcher.RemoveConnectionBanMeta{
			User:   "someUser",
			Reason: "someReason",
		},
	})
	assert.Error(t, err)
	assert.Equal(t, "forced remove connection ban error", err.Error())

	assemblyStatusOutput, err := dsp.AssemblyStatus(context.Background())
	require.NoError(t, err)
	assert.Equal(t, dispatcher.AssemblyStateInProgress, assemblyStatusOutput.State)
	require.NotNil(t, assemblyStatusOutput.EstimatedTimeRemaining)
	assert.Equal(t, 100, *assemblyStatusOutput.EstimatedTimeRemaining)

	reviewCasesOutput, err := dsp.ReviewCases(context.Background(), &dispatcher.ReviewCasesInput{})
	require.NoError(t, err)
	require.Len(t, reviewCasesOutput.Cases, 1)
	assert.Equal(t, &testReviewCase, reviewCasesOutput.Cases[0])
	require.NotNil(t, reviewCasesOutput.NextCursor)
	assert.Equal(t, "someCursor", *reviewCasesOutput.NextCursor)
	assert.Equal(t, 3, reviewCasesOutput.OpenCount)

	claimOutput, err := dsp.ClaimReviewCase(context.Background(), &dispatcher.ClaimReviewCaseInput{
		ID:    "someCase",
		Actor: testReviewActor,
	})
	require.NoError(t, err)
	assert.Equal(t, &testReviewCase, claimOutput.Case)
	assert.Equal(t, "someEntity", claimOutput.EntityID)
	assert.Equal(t, "review:someUser:1", claimOutput.Lock)
	assert.False(t, claimOutput.Stale)

	releaseOutput, err := dsp.ReleaseReviewCase(context.Background(), &dispatcher.ReleaseReviewCaseInput{
		ID:       "someCase",
		EntityID: "someEntity",
		Lock:     "review:someUser:1",
	})
	require.NoError(t, err)
	assert.True(t, releaseOutput.Released)

	resolveOutput, err := dsp.ResolveReviewCase(context.Background(), &dispatcher.ResolveReviewCaseInput{
		ID: "someCase",
		Verdicts: []dispatcher.ReviewCandidateVerdict{
			{
				EntityID: "someEntity",
				Keep:     false,
			},
		},
		Actor:  testReviewActor,
		Reason: "someReason",
	})
	require.NoError(t, err)
	assert.True(t, resolveOutput.Triggered)
}

type testDispatcher struct {
	deadlineExists  bool
	disassembleLock string
}

var testEntity = api.Entity{
	ID: "abcd",
	Records: []*api.Record{
		{
			ID: "12345",
			Data: map[string]any{
				"foo": "bar",
			},
		},
	},
	Edges: api.Edges{
		"12345:12345:STATIC",
	},
	Duplicates: api.Duplicates{
		"12345": []string{
			"12345",
			"duplicateID",
		},
	},
	Hits: api.Hits{
		"12345": []string{"someRuleName"},
	},
}

var testReviewActor = dispatcher.ReviewActor{
	Kind: dispatcher.ReviewActorHuman,
	ID:   "someUser",
	Name: "Some User",
}

var testReviewCase = dispatcher.ReviewCase{
	ID:        "someCase",
	Status:    dispatcher.ReviewCaseClaimed,
	CreatedAt: time.Date(2026, 8, 19, 10, 0, 0, 0, time.UTC),
	RecordIDs: []string{"12345"},
	EntityID:  "someEntity",
	Candidates: []dispatcher.ReviewCandidate{
		{
			EntityID: "someOtherEntity",
			Links: []dispatcher.ReviewLink{
				{
					A: "12345",
					B: "67890",
					Rules: []dispatcher.ReviewRule{
						{
							ID:    "someRuleName",
							Score: 100,
						},
					},
				},
			},
		},
	},
	ClaimedBy: &testReviewActor,
	ClaimedAt: &testReviewClaimedAt,
}

var testReviewClaimedAt = time.Date(2026, 8, 19, 11, 0, 0, 0, time.UTC)

func (d *testDispatcher) Entity(ctx context.Context, _ *dispatcher.EntityInput) (*dispatcher.EntityOutput, error) {
	_, d.deadlineExists = ctx.Deadline()
	return &dispatcher.EntityOutput{
		Entity: &testEntity,
	}, nil
}

func (d *testDispatcher) EntityByRecord(ctx context.Context, _ *dispatcher.EntityByRecordInput) (*dispatcher.EntityOutput, error) {
	_, d.deadlineExists = ctx.Deadline()
	return &dispatcher.EntityOutput{
		Entity: &testEntity,
	}, nil
}

func (d *testDispatcher) Search(_ context.Context, _ *dispatcher.SearchInput) (*dispatcher.SearchOutput, error) {
	return &dispatcher.SearchOutput{
		Entities: []*api.Entity{
			&testEntity,
		},
	}, nil
}

func (d *testDispatcher) Submit(_ context.Context, _ *dispatcher.SubmitInput) (*dispatcher.SubmitOutput, error) {
	return &dispatcher.SubmitOutput{
		RecordsAdded: 1,
	}, nil
}

func (d *testDispatcher) SubmitWithPreview(ctx context.Context, input *dispatcher.SubmitWithPreviewInput) (*dispatcher.SubmitWithPreviewOutput, error) {
	return &dispatcher.SubmitWithPreviewOutput{
		Entities: []*api.Entity{
			&testEntity,
		},
	}, nil
}

func (d *testDispatcher) Disassemble(_ context.Context, input *dispatcher.DisassembleInput) (*dispatcher.DisassembleOutput, error) {
	d.disassembleLock = input.Lock
	return &dispatcher.DisassembleOutput{
		Triggered: true,
	}, nil
}

func (d *testDispatcher) RemoveConnectionBan(_ context.Context, _ *dispatcher.RemoveConnectionBanInput) error {
	return fmt.Errorf("forced remove connection ban error")
}

func (d *testDispatcher) AssemblyStatus(_ context.Context) (*dispatcher.AssemblyStatusOutput, error) {
	t := 100
	return &dispatcher.AssemblyStatusOutput{
		State:                  dispatcher.AssemblyStateInProgress,
		EstimatedTimeRemaining: &t,
	}, nil
}

func (d *testDispatcher) ReviewCases(_ context.Context, _ *dispatcher.ReviewCasesInput) (*dispatcher.ReviewCasesOutput, error) {
	cursor := "someCursor"
	return &dispatcher.ReviewCasesOutput{
		Cases:      []*dispatcher.ReviewCase{&testReviewCase},
		NextCursor: &cursor,
		OpenCount:  3,
	}, nil
}

func (d *testDispatcher) ClaimReviewCase(_ context.Context, _ *dispatcher.ClaimReviewCaseInput) (*dispatcher.ClaimReviewCaseOutput, error) {
	return &dispatcher.ClaimReviewCaseOutput{
		Case:     &testReviewCase,
		EntityID: "someEntity",
		Lock:     "review:someUser:1",
	}, nil
}

func (d *testDispatcher) ReleaseReviewCase(_ context.Context, _ *dispatcher.ReleaseReviewCaseInput) (*dispatcher.ReleaseReviewCaseOutput, error) {
	return &dispatcher.ReleaseReviewCaseOutput{
		Released: true,
	}, nil
}

func (d *testDispatcher) ResolveReviewCase(_ context.Context, _ *dispatcher.ResolveReviewCaseInput) (*dispatcher.ResolveReviewCaseOutput, error) {
	return &dispatcher.ResolveReviewCaseOutput{
		Triggered: true,
	}, nil
}
