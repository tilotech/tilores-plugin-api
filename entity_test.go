package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	api "github.com/tilotech/tilores-plugin-api"
)

func TestParseEdge(t *testing.T) {
	edge := "id1:0:id2:9:R1"
	a, b, rule := api.ParseEdge(edge)
	assert.Equal(t, "id1:0", a)
	assert.Equal(t, "id2:9", b)
	assert.Equal(t, "R1", rule)
}

func TestParseDuplicateKey(t *testing.T) {
	duplicate := ":id:9"
	id, group := api.ParseDuplicateKey(duplicate)
	assert.Equal(t, "id:9", id)
	assert.Equal(t, "", group)

	duplicate = "default:id:9"
	id, group = api.ParseDuplicateKey(duplicate)
	assert.Equal(t, "id:9", id)
	assert.Equal(t, "default", group)
}

func TestParseRecordID(t *testing.T) {
	rid := "foo:9"
	id, version := api.ParseRecordID(rid)
	assert.Equal(t, "foo", id)
	assert.Equal(t, 9, version)
}

func TestNewEdge(t *testing.T) {
	actual := api.NewEdge("foo:1", "bar:2", "R1")
	assert.Equal(t, "foo:1:bar:2:R1", actual)

	actual = api.NewEdgeWithVersions("foo", 1, "bar", 2, "R2")
	assert.Equal(t, "foo:1:bar:2:R2", actual)
}

func TestNewDuplicateKey(t *testing.T) {
	actual := api.NewDuplicateKey("foo:1", "")
	assert.Equal(t, ":foo:1", actual)

	actual = api.NewDuplicateKey("foo:1", "default")
	assert.Equal(t, "default:foo:1", actual)
}

func TestNewRecordID(t *testing.T) {
	actual := api.NewRecordID("foo", 9)
	assert.Equal(t, "foo:9", actual)
}
