//go:build integration

package greeter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGreeter(t *testing.T) {
	t.Run("should return greetings", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/greetings", httpPort))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var res []string
		err = json.NewDecoder(resp.Body).Decode(&res)
		require.NoError(t, err)

		assert.Len(t, res, 2)
	})
}
