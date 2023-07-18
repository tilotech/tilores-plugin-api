package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	api "github.com/tilotech/tilores-plugin-api"
)

func TestIDWithVersion(t *testing.T) {
	record := &api.Record{ID: "foo"}
	assert.Equal(t, "foo", record.IDWithVersion())

	record = &api.Record{ID: "foo", Meta: &api.RecordMeta{Version: 0}}
	assert.Equal(t, "foo", record.IDWithVersion())

	record = &api.Record{ID: "foo", Meta: &api.RecordMeta{Version: 1}}
	assert.Equal(t, "foo:1", record.IDWithVersion())
}
