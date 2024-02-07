package tests

import (
	"testing"

	"github.com/royalcat/query"
	"github.com/stretchr/testify/require"
)

type id int64

type model struct {
	ID     id          `json:"id"`
	Name   string      `json:"name"`
	Nested nestedModel `json:"nested"`
}

type nestedModel struct {
	Based bool `json:"based"`
}

func TestParseStringFilter(t *testing.T) {
	require := require.New(t)
	sf := map[string]string{
		"id{in}":           "69,420",
		"name{substr}":     "Primagen",
		"nested.based{eq}": "true",
	}
	f, err := query.ParseStringFilter[model](sf)
	require.NoError(err)
	require.ElementsMatch(query.Filter{
		{Field: "id", Op: query.OperatorIn, Value: []id{69, 420}},
		{Field: "name", Op: query.OperatorSubString, Value: "Primagen"},
		{Field: "nested.based", Op: query.OperatorEqual, Value: true},
	}, f)
}
