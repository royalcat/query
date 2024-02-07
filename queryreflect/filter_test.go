package queryreflect_test

import (
	"testing"

	"github.com/royalcat/query"
	"github.com/royalcat/query/queryreflect"
	"github.com/stretchr/testify/require"
)

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
		{Field: "names", Value: 1},
	}

	out, err := queryreflect.ApplyFilter(f, data)
	require.NoError(err)
	require.Equal([]testStruct{
		{Names: []int{1, 2, 3}},
		{Names: []int{1, 4}},
	}, out)
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
		{Field: "name", Op: query.OperatorSubString, Value: "a"},
	}
	out, err := queryreflect.ApplyFilter(f, data)
	require.NoError(err)
	require.Equal([]testStruct{
		{ID: 4, Name: "aac"},
		{ID: 2, Name: "aab"},
		{ID: 1, Name: "aaa"},
	}, out)
}

func BenchmarkApplyFilter(b *testing.B) {
	type testStruct struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	data := []testStruct{
		{ID: 4, Name: "aac"},
		{ID: 3, Name: "bbb"},
		{ID: 2, Name: "aab"},
		{ID: 1, Name: "aaa"},
	}
	f := query.Filter{
		{Field: "name", Op: query.OperatorSubString, Value: "a"},
	}
	// require.Equal([]testStruct{
	// 	{ID: 4, Name: "aac"},
	// 	{ID: 2, Name: "aab"},
	// 	{ID: 1, Name: "aaa"},
	// }, out)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := queryreflect.ApplyFilter(f, data)
		if err != nil {
			b.Fatalf("error applying filter: %s", err.Error())
		}
	}

}
