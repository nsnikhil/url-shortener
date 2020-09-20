package router_test

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
	"urlshortner/pkg/elongator"
	router2 "urlshortner/pkg/http/router"
	"urlshortner/pkg/reporters"
	"urlshortner/pkg/shortener"
)

func TestRouter(t *testing.T) {
	r := router2.NewRouter(
		zap.NewNop(),
		&newrelic.Application{},
		&reporters.MockPrometheus{},
		&shortener.MockShortener{},
		&elongator.MockElongator{},
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
