package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/davidklassen/testing-workshop/pkg/greetings"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRouter(t *testing.T) {
	t.Run("should return greetings", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := NewMockGreetingsRepository(ctrl)
		greeter := NewGreeter(repo, 1, time.Millisecond, 2)
		router := NewRouter(greeter)

		server := httptest.NewServer(router)

		repo.EXPECT().FindAll().Times(1).Return([]*greetings.Greeting{{}}, nil)

		resp, err := http.Get(fmt.Sprintf("%s/greetings", server.URL))
		require.NoError(t, err)

		res := make([]string, 0)

		err = json.NewDecoder(resp.Body).Decode(&res)
		require.NoError(t, err)
		assert.Len(t, res, 1)
	})
}
