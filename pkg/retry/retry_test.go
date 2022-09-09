package retry_test

import (
	"errors"
	"testing"
	"time"

	"github.com/davidklassen/testing-workshop/pkg/retry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetry(t *testing.T) {
	t.Run("should retry max attempts times", func(t *testing.T) {
		attempts := 0
		err := retry.Retry(5, time.Millisecond, 1, func() error {
			attempts++
			return errors.New("error")
		})
		require.Error(t, err)
		assert.Equal(t, 5, attempts)
	})

	t.Run("should return original error", func(t *testing.T) {
		origErr := errors.New("error")
		err := retry.Retry(1, time.Millisecond, 1, func() error {
			return origErr
		})
		assert.Equal(t, err, origErr)
	})

	t.Run("should stop on success", func(t *testing.T) {
		attempts := 0
		err := retry.Retry(5, time.Millisecond, 1, func() error {
			attempts++
			if attempts < 3 {
				return errors.New("error")
			}
			return nil
		})
		require.NoError(t, err)
		assert.Equal(t, 3, attempts)
	})
}
