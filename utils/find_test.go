package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	items := []string{"foo", "bar"}
	t.Run("with positive found", func(t *testing.T) {
		result := Find(items, func(item string) bool { return item == "bar" })
		assert.Equal(t, "bar", result)
	})

	t.Run("with negative found", func(t *testing.T) {
		result := Find(items, func(item string) bool { return item == "foobar" })
		assert.Equal(t, "", result)
	})

	t.Run("with complex types", func(t *testing.T) {
		type Person struct {
			Name string
		}
		persons := []Person{{Name: "Paul"}, {Name: "Simone"}}
		result := Find(persons, func(person Person) bool { return person.Name == "Simone" })
		assert.Equal(t, persons[1], result)
	})
}
