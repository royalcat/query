package query_test

import (
	"testing"

	"github.com/royalcat/query"
	"github.com/stretchr/testify/require"
)

func TestFilterFields(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    query.Filter
		expected query.Fields
	}{
		{
			name:     "When Filter is empty",
			input:    query.Filter{},
			expected: query.Fields{},
		},
		{
			name:     "When Filter has one field",
			input:    query.Filter{"field1": "value1"},
			expected: query.Fields{"field1"},
		},
		{
			name:     "When Filter has multiple fields",
			input:    query.Filter{"field1": "value1", "field2": "value2"},
			expected: query.Fields{"field1", "field2"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, tc.input.Fields())
		})
	}
}
