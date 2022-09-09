package greetings_test

import (
	"testing"

	"github.com/davidklassen/testing-workshop/pkg/greetings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGreeting(t *testing.T) {
	t.Run("should create a greeting", func(t *testing.T) {
		_, err := greetings.NewGreeting("Hello, {{.}}!")
		require.NoError(t, err)
	})

	t.Run("should fail with incorrect template", func(t *testing.T) {
		_, err := greetings.NewGreeting("Hello, {{undefinedFunc}}!")
		assert.Error(t, err)
	})
}

func TestGreeting_Greet(t *testing.T) {
	t.Run("should greet", func(t *testing.T) {
		greeting, err := greetings.NewGreeting("Hello, {{.}}!")
		require.NoError(t, err)
		res, err := greeting.Greet("World")
		require.NoError(t, err)
		assert.Equal(t, "Hello, World!", res)
	})

	t.Run("should fail with invalid input", func(t *testing.T) {
		greeting, err := greetings.NewGreeting("Hello, {{.}}!")
		require.NoError(t, err)
		_, err = greeting.Greet("")
		assert.Error(t, err)
	})
}
