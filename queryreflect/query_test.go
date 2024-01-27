package queryreflect_test

import (
	"testing"

	"github.com/royalcat/query"
	"github.com/royalcat/query/queryreflect"

	"github.com/stretchr/testify/require"
)

func TestApplyQuery(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	require := require.New(t)
	data := []testStruct{
		{ID: 5, Name: "bbc"},
		{ID: 4, Name: "aac"},
		{ID: 3, Name: "bbb"},
		{ID: 2, Name: "aab"},
		{ID: 1, Name: "aaa"},
	}
	q := query.Query{
		Filter: query.Filter{
			"name{substr}": "a",
		},
		Sort: query.Sort{
			{Key: "id", Order: query.ASC},
		},
		Pagination: query.Pagination{
			Offset: 1,
			Limit:  2,
		},
	}
	out, err := queryreflect.ApplyQuery(q, data)
	require.NoError(err)
	require.Equal([]testStruct{
		{ID: 2, Name: "aab"},
		{ID: 4, Name: "aac"},
	}, out)
}

func TestApplyQueryArrayIndex(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Names []int `json:"names"`
	}

	require := require.New(t)
	data := []testStruct{
		{Names: []int{3}},
		{Names: []int{}},
		{Names: []int{1}},
		{Names: []int{2}},
		{Names: []int{4}},
	}
	f := query.Query{
		Filter: query.Filter{
			"names.0{gt}": "1",
		},
		Sort: query.Sort{
			{Key: "names.0", Order: query.ASC},
		},
	}

	out, err := queryreflect.ApplyQuery(f, data)
	require.NoError(err)
	require.Equal([]testStruct{
		{Names: []int{2}},
		{Names: []int{3}},
		{Names: []int{4}},
	}, out)
}
