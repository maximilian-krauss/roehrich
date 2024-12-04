package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsAll(t *testing.T) {
	type TestCase struct {
		name     string
		subset   []string
		superset []string
		expected bool
	}

	testCases := []TestCase{
		{name: "Partial matching", subset: []string{"a", "b"}, superset: []string{"a", "b", "c"}, expected: true},
		{name: "Complete matching", subset: []string{"a", "b", "c"}, superset: []string{"a", "b", "c"}, expected: true},
		{name: "Partial miss", subset: []string{"a", "b", "c", "d"}, superset: []string{"a", "b", "c"}, expected: false},
		{name: "Complete miss", subset: []string{"x", "y", "z"}, superset: []string{"a", "b", "c"}, expected: false},
		{name: "Empty subset", subset: []string{}, superset: []string{"a", "b", "c"}, expected: true},
		{name: "Empty superset", subset: []string{"x", "y", "z"}, superset: []string{}, expected: false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			println(testCase.name)
			result := ContainsAll(testCase.subset, testCase.superset)
			assert.Equal(t, testCase.expected, result)
		})
	}
}
