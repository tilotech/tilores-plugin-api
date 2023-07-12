package api_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
	api "github.com/tilotech/tilores-plugin-api"
)

func TestHitsIDs(t *testing.T) {
	hits := api.Hits{
		"id1":   []string{"R1", "R2"},
		"id2:0": []string{"R1"},
		"id3:1": []string{"R1"},
	}

	actual := hits.IDs()
	assert.ElementsMatch(t, []string{"id1", "id2", "id3"}, actual)
}

func TestParseEdge(t *testing.T) {
	edge := "id1:0:id2:9:R1"
	a, b, rule := api.ParseEdge(edge)
	assert.Equal(t, "id1", a)
	assert.Equal(t, "id2:9", b)
	assert.Equal(t, "R1", rule)

	edge = "id1:id2:R1"
	a, b, rule = api.ParseEdge(edge)
	assert.Equal(t, "id1", a)
	assert.Equal(t, "id2", b)
	assert.Equal(t, "R1", rule)
}

func TestParseDuplicateKey(t *testing.T) {
	duplicate := ":id:9"
	id, group := api.ParseDuplicateKey(duplicate)
	assert.Equal(t, "id:9", id)
	assert.Equal(t, "", group)

	duplicate = ":id:0"
	id, group = api.ParseDuplicateKey(duplicate)
	assert.Equal(t, "id", id)
	assert.Equal(t, "", group)

	duplicate = "default:id:9"
	id, group = api.ParseDuplicateKey(duplicate)
	assert.Equal(t, "id:9", id)
	assert.Equal(t, "default", group)

	duplicate = "default:id"
	id, group = api.ParseDuplicateKey(duplicate)
	assert.Equal(t, "id", id)
	assert.Equal(t, "default", group)

	duplicate = "id"
	id, group = api.ParseDuplicateKey(duplicate)
	assert.Equal(t, "id", id)
	assert.Equal(t, "", group)
}

func TestParseRecordID(t *testing.T) {
	rid := "foo:9"
	id, version := api.ParseRecordID(rid)
	assert.Equal(t, "foo", id)
	assert.Equal(t, 9, version)

	rid = "foo:0"
	id, version = api.ParseRecordID(rid)
	assert.Equal(t, "foo", id)
	assert.Equal(t, 0, version)

	rid = "foo"
	id, version = api.ParseRecordID(rid)
	assert.Equal(t, "foo", id)
	assert.Equal(t, 0, version)
}

func TestParseRecordIDWithOptionalVersion(t *testing.T) {
	rid := "foo"
	id, version := api.ParseRecordIDWithOptionalVersion(rid)
	assert.Equal(t, "foo", id)
	assert.Nil(t, version)

	rid = "foo:0"
	id, version = api.ParseRecordIDWithOptionalVersion(rid)
	assert.Equal(t, "foo", id)
	require.NotNil(t, version)
	assert.Equal(t, 0, *version)

	rid = "foo:1"
	id, version = api.ParseRecordIDWithOptionalVersion(rid)
	assert.Equal(t, "foo", id)
	require.NotNil(t, version)
	assert.Equal(t, 1, *version)
}

func TestNewEdge(t *testing.T) {
	actual := api.NewEdge("foo:1", "bar:2", "R1")
	assert.Equal(t, "foo:1:bar:2:R1", actual)

	actual = api.NewEdge("foo", "bar:2", "R1")
	assert.Equal(t, "foo:0:bar:2:R1", actual)

	actual = api.NewEdge("foo", "bar", "R1")
	assert.Equal(t, "foo:0:bar:0:R1", actual)

	actual = api.NewEdgeWithVersions("foo", 1, "bar", 2, "R2")
	assert.Equal(t, "foo:1:bar:2:R2", actual)
}

func TestNewDuplicateKey(t *testing.T) {
	actual := api.NewDuplicateKey("foo:1", "")
	assert.Equal(t, ":foo:1", actual)

	actual = api.NewDuplicateKey("foo", "")
	assert.Equal(t, ":foo:0", actual)

	actual = api.NewDuplicateKey("foo:1", "default")
	assert.Equal(t, "default:foo:1", actual)

	actual = api.NewDuplicateKey("foo", "default")
	assert.Equal(t, "default:foo:0", actual)
}

func TestNewRecordID(t *testing.T) {
	actual := api.NewRecordID("foo", 9)
	assert.Equal(t, "foo:9", actual)

	actual = api.NewRecordID("foo", 0)
	assert.Equal(t, "foo", actual)
}

func TestNewRecordIDWithVersion(t *testing.T) {
	actual := api.NewRecordIDWithVersion("foo", 9)
	assert.Equal(t, "foo:9", actual)

	actual = api.NewRecordIDWithVersion("foo", 0)
	assert.Equal(t, "foo:0", actual)
}
