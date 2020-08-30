package router_test

import (
	"fmt"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gopkg.in/alexcesaro/statsd.v2"
	"net/http"
	"net/http/httptest"
	"testing"
	"urlshortner/pkg/shortener/router"
	"urlshortner/pkg/shortener/service"
)

func TestRouter(t *testing.T) {
	r := router.NewRouter(
		zap.NewNop(),
		&newrelic.Application{},
		&statsd.Client{},
		service.NewService(&service.MockShortenerService{}),
	)

	testCases := []struct {
		name         string
		actualResult func() int
	}{
		{
			name: "test ping route",
			actualResult: func() int {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest(http.MethodGet, "/ping", nil)
				require.NoError(t, err)

				r.ServeHTTP(resp, req)

				return resp.Code
			},
		},
		{
			name: "test shorten route",
			actualResult: func() int {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest(http.MethodPost, "/shorten", nil)
				require.NoError(t, err)

				r.ServeHTTP(resp, req)

				fmt.Println(resp)
				return resp.Code
			},
		},
		{
			name: "test redirect route",
			actualResult: func() int {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest(http.MethodGet, "/IxGtJ", nil)
				require.NoError(t, err)

				r.ServeHTTP(resp, req)

				fmt.Println(resp)
				return resp.Code
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.NotEqual(t, http.StatusNotFound, testCase.actualResult())
		})
	}
}
