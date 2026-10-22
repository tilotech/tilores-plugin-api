package dispatcher_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "github.com/tilotech/tilores-plugin-api"
	"github.com/tilotech/tilores-plugin-api/dispatcher"
)

func TestUnmarshalEvent(t *testing.T) {
	cases := map[string]struct {
		input       string
		expected    *dispatcher.AssembleEvent
		expectError bool
	}{
		"standard assemble": {
			input: `
				{
					"type": "ASSEMBLE",
					"payload": [
						{
							"id": "foo",
							"data": {
								"key": "value"
							}
						}
					]
				}`,
			expected: &dispatcher.AssembleEvent{
				Type: "ASSEMBLE",
				Payload: []*api.Record{
					{
						ID: "foo",
						Data: map[string]any{
							"key": "value",
						},
					},
				},
			},
		},
		"standard disassemble": {
			input: `
				{
					"type": "DISASSEMBLE",
					"payload": {
						"edges": [
							{
								"a": "foo-1",
								"b": "foo-2"
							}
						]
					}
				}`,
			expected: &dispatcher.AssembleEvent{
				Type: "DISASSEMBLE",
				Payload: &dispatcher.DisassembleInput{
					Edges: []dispatcher.DisassembleEdge{
						{
							A: "foo-1",
							B: "foo-2",
						},
					},
				},
			},
		},
		"plain outdated assemble": {
			input: `
				[
					{
						"id": "foo",
						"data": {
							"key": "value"
						}
					}
				]`,
			expected: &dispatcher.AssembleEvent{
				Type: "ASSEMBLE",
				Payload: []*api.Record{
					{
						ID: "foo",
						Data: map[string]any{
							"key": "value",
						},
					},
				},
			},
		},
		"assemble with state rectification flag": {
			input: `
				{
					"type": "ASSEMBLE",
					"payload": [
						{
							"id": "foo",
							"data": {
								"key": "value"
							}
						}
					],
					"flags": ["STATE_RECTIFICATION"]
				}`,
			expected: &dispatcher.AssembleEvent{
				Type: "ASSEMBLE",
				Payload: []*api.Record{
					{
						ID: "foo",
						Data: map[string]any{
							"key": "value",
						},
					},
				},
				Flags: []dispatcher.AssembleEventFlag{dispatcher.StateRectification},
			},
		},
		"review assemble": {
			input: `
				{
					"type": "REVIEW_ASSEMBLE",
					"payload": {
						"caseID": "case-1",
						"records": [
							{
								"id": "foo",
								"data": {
									"key": "value"
								},
								"meta": {
									"version": 2
								}
							}
						],
						"verdicts": [
							{
								"a": "foo:2",
								"b": "bar:0",
								"keep": true
							}
						],
						"lock": "review:actor-1:1700000000",
						"lockedEntityIDs": ["entity-1", "entity-2"]
					}
				}`,
			expected: &dispatcher.AssembleEvent{
				Type: "REVIEW_ASSEMBLE",
				Payload: &dispatcher.ReviewAssembleInput{
					CaseID: "case-1",
					Records: []*api.Record{
						{
							ID: "foo",
							Data: map[string]any{
								"key": "value",
							},
							Meta: &api.RecordMeta{
								Version: 2,
							},
						},
					},
					Verdicts: []dispatcher.ReviewLinkVerdict{
						{
							A:    "foo:2",
							B:    "bar:0",
							Keep: true,
						},
					},
					Lock:            "review:actor-1:1700000000",
					LockedEntityIDs: []string{"entity-1", "entity-2"},
				},
			},
		},
		"invalid review assemble payload": {
			input:       `{"type": "REVIEW_ASSEMBLE", "payload": []}`,
			expectError: true,
		},
		"invalid json": {
			input:       `{invalid json`,
			expectError: true,
		},
		"invalid payload": {
			input:       `{"type": "ASSEMBLE", "payload": {}}`,
			expectError: true,
		},
		"invalid type": {
			input:       `{"type": "UNKNOWN"}`,
			expectError: true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			request := &dispatcher.AssembleEvent{}
			err := json.Unmarshal([]byte(c.input), request)
			if c.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, c.expected, request)
		})
	}
}
