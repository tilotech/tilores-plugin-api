package dispatcher_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "github.com/tilotech/tilores-plugin-api"
	"github.com/tilotech/tilores-plugin-api/dispatcher"
	"testing"
)

func TestPlugin(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reattachConfigCh := make(chan *plugin.ReattachConfig, 1)
	go providePluginServer(ctx, reattachConfigCh)

	dsp, err := createPluginClient(reattachConfigCh)
	require.NoError(t, err)

	entity, err := dsp.Entity(context.Background(), "abcd")
	assert.NoError(t, err)
	assert.NotNil(t, entity)
	assert.Equal(t, 1, len(entity.Records))
	assert.Equal(t, "bar", entity.Records[0].Data["foo"])
	assert.Equal(t, 1, len(entity.Edges))
	assert.Equal(t, 1, len(entity.Duplicates))

	parameters := &api.SearchParameters{
		"foo": "bar",
	}
	entities, err := dsp.Search(context.Background(), parameters)
	assert.NoError(t, err)
	assert.NotNil(t, entities)
	assert.Equal(t, 1, len(entities))

	submitOutput, err := dsp.Submit(context.Background(), &dispatcher.SubmitInput{
		Records: []*api.Record{
			{
				ID: "12345", Data: map[string]interface{}{
					"foo": "bar"},
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, submitOutput.RecordsAdded)

	disassembleOutput, err := dsp.Disassemble(context.Background(), &dispatcher.DisassembleInput{
		Reference: "123123",
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
		Meta: dispatcher.DisassembleMeta{
			User:   "someUser",
			Reason: "someReason",
		},
		Timeout: nil,
	})
	assert.NoError(t, err)
	assert.Equal(t, []string{"abcd"}, disassembleOutput.EntityIDs)

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

func providePluginServer(ctx context.Context, reattachConfigCh chan<- *plugin.ReattachConfig) {
	var pluginMap = map[string]plugin.Plugin{
		"dispatcher": &dispatcher.Plugin{
			Impl: &testDispatcher{},
		},
	}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: dispatcher.Handshake,
		Plugins:         pluginMap,
		Test: &plugin.ServeTestConfig{
			Context:          ctx,
			ReattachConfigCh: reattachConfigCh,
		},
	})
}

func createPluginClient(reattachConfigCh chan *plugin.ReattachConfig) (dispatcher.Dispatcher, error) {
	reattachConfig := <-reattachConfigCh
	if reattachConfig == nil {
		return nil, fmt.Errorf("expected reattach config, but received none")
	}

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: dispatcher.Handshake,
		Plugins: map[string]plugin.Plugin{
			"dispatcher": &dispatcher.Plugin{},
		},
		Reattach: reattachConfig,
	})

	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}

	raw, err := rpcClient.Dispense("dispatcher")
	if err != nil {
		return nil, err
	}

	impl, ok := raw.(dispatcher.Dispatcher)
	if !ok {
		return nil, fmt.Errorf("not a dispatcher plugin: %T", raw)
	}

	return impl, nil
}

type testDispatcher struct{}

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
}

func (d *testDispatcher) Entity(_ context.Context, _ string) (*api.Entity, error) {
	return &testEntity, nil
}

func (d *testDispatcher) Search(_ context.Context, _ *api.SearchParameters) ([]*api.Entity, error) {
	return []*api.Entity{
		&testEntity,
	}, nil
}

func (d *testDispatcher) Submit(_ context.Context, _ *dispatcher.SubmitInput) (*dispatcher.SubmitOutput, error) {
	return &dispatcher.SubmitOutput{
		RecordsAdded: 1,
	}, nil
}

func (d *testDispatcher) Disassemble(_ context.Context, _ *dispatcher.DisassembleInput) (*dispatcher.DisassembleOutput, error) {
	return &dispatcher.DisassembleOutput{
		DeletedEdges:   1,
		DeletedRecords: 1,
		EntityIDs:      []string{"abcd"},
	}, nil
}

func (d *testDispatcher) RemoveConnectionBan(_ context.Context, _ *dispatcher.RemoveConnectionBanInput) error {
	return fmt.Errorf("forced remove connection ban error")
}
