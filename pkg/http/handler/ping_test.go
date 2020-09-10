package handler_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
	"urlshortner/pkg/http/handler"
	"urlshortner/pkg/reporters"
)

func TestPingHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/ping", nil)
	require.NoError(t, err)

	mockStatsD := &reporters.MockStatsDClient{}
	mockStatsD.On("ReportAttempt", "ping")
	mockStatsD.On("ReportSuccess", "ping")

	handler.PingHandler(zap.NewNop(), mockStatsD)(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
