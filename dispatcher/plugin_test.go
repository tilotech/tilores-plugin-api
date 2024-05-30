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
				Data: map[string]interface{}{
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
				Data: map[string]interface{}{
					"foo": "bar",
				},
			},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, submitWithPreviewOutput)
	assert.Len(t, submitWithPreviewOutput.Preview.Entities, 1)

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
	})
	assert.NoError(t, err)
	assert.True(t, disassembleOutput.Triggered)

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
}

type testDispatcher struct {
	deadlineExists bool
}

var testEntity = api.Entity{
	ID: "abcd",
	Records: []*api.Record{
		{
			ID: "12345",
			Data: map[string]interface{}{
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
		Preview: dispatcher.SubmissionPreview{
			Entities: []*api.Entity{
				&testEntity,
			},
			NewRecords:     nil,
			UpdatedRecords: []string{"12345"},
			IgnoredRecords: nil,
		},
	}, nil
}

func (d *testDispatcher) Disassemble(_ context.Context, _ *dispatcher.DisassembleInput) (*dispatcher.DisassembleOutput, error) {
	return &dispatcher.DisassembleOutput{
		Triggered: true,
	}, nil
}

func (d *testDispatcher) RemoveConnectionBan(_ context.Context, _ *dispatcher.RemoveConnectionBanInput) error {
	return fmt.Errorf("forced remove connection ban error")
}
