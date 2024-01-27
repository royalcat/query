package queryreflect_test

import (
	"testing"

	"github.com/royalcat/query"
	"github.com/royalcat/query/queryreflect"
	"github.com/stretchr/testify/require"
)

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
