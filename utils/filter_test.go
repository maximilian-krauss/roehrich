package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilter(t *testing.T) {
	t.Run("test with strings", func(t *testing.T) {
		input := []string{"one", "two", "three"}
		expected := []string{"three"}
		output := Filter(input, func(s string) bool { return len(s) > 3 })

		assert.Equal(t, expected, output)
	})

	t.Run("test with numbers", func(t *testing.T) {
		input := []int{1, 2, 3}
		expected := []int{3}
		output := Filter(input, func(i int) bool { return i >= 3 })

		assert.Equal(t, expected, output)
	})

	t.Run("test with empty output", func(t *testing.T) {
		input := []int{1, 2, 3}
		var expected []int
		output := Filter(input, func(i int) bool { return i > 3 })

		assert.Equal(t, expected, output)
	})
}
