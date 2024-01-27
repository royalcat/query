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

func TestApplyFilterArrayAny(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Names []int `json:"names"`
	}

	require := require.New(t)
	data := []testStruct{
		{Names: []int{}},
		{Names: []int{1, 2, 3}},
		{Names: []int{2, 3, 4}},
		{Names: []int{1, 4}},
	}
	f := query.Filter{
		"names": "1",
	}

	out, err := queryreflect.ApplyFilter(f, data)
	require.NoError(err)
	require.Equal([]testStruct{
		{Names: []int{1, 2, 3}},
		{Names: []int{1, 4}},
	}, out)
}

func TestApplySort(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	require := require.New(t)
	{
		data := []testStruct{
			{ID: 5, Name: "e"},
			{ID: 4, Name: "d"},
			{ID: 3, Name: "c"},
			{ID: 2, Name: "b"},
			{ID: 1, Name: "a"},
		}
		s := query.Sort{
			{Key: "name", Order: query.ASC},
		}

		out, err := queryreflect.ApplySort(s, data)
		require.NoError(err)
		require.Equal([]testStruct{
			{ID: 1, Name: "a"},
			{ID: 2, Name: "b"},
			{ID: 3, Name: "c"},
			{ID: 4, Name: "d"},
			{ID: 5, Name: "e"},
		}, out)
	}
	{
		data := []testStruct{
			{ID: 1, Name: "test3"},
			{ID: 3, Name: "test1"},
			{ID: 5, Name: "apollo"},
			{ID: 2, Name: "test2"},
			{ID: 4, Name: "royalcat"},
		}
		s := query.Sort{
			{Key: "name", Order: query.ASC},
		}
		out, err := queryreflect.ApplySort(s, data)
		require.NoError(err)
		require.Equal([]testStruct{
			{ID: 5, Name: "apollo"},
			{ID: 4, Name: "royalcat"},
			{ID: 3, Name: "test1"},
			{ID: 2, Name: "test2"},
			{ID: 1, Name: "test3"},
		}, out)
	}
	{
		data := []testStruct{
			{ID: 1, Name: "test3"},
			{ID: 1, Name: "test1"},
			{ID: 1, Name: "apollo"},
			{ID: 1, Name: "test2"},
			{ID: 1, Name: "royalcat"},
		}
		s := query.Sort{
			{Key: "name", Order: query.DESC},
		}
		out, err := queryreflect.ApplySort(s, data)
		require.NoError(err)
		require.Equal([]testStruct{
			{ID: 1, Name: "test3"},
			{ID: 1, Name: "test2"},
			{ID: 1, Name: "test1"},
			{ID: 1, Name: "royalcat"},
			{ID: 1, Name: "apollo"},
		}, out)
	}
}
func TestApplyFilter(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	require := require.New(t)
	data := []testStruct{
		{ID: 4, Name: "aac"},
		{ID: 3, Name: "bbb"},
		{ID: 2, Name: "aab"},
		{ID: 1, Name: "aaa"},
	}
	f := query.Filter{
		"name{substr}": "a",
	}
	out, err := queryreflect.ApplyFilter(f, data)
	require.NoError(err)
	require.Equal([]testStruct{
		{ID: 4, Name: "aac"},
		{ID: 2, Name: "aab"},
		{ID: 1, Name: "aaa"},
	}, out)
}
