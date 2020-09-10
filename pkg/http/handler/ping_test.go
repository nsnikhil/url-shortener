package handler_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gopkg.in/alexcesaro/statsd.v2"
	"net/http"
	"net/http/httptest"
	"testing"
	"urlshortner/pkg/http/handler"
)

func TestPingHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/ping", nil)
	require.NoError(t, err)

	handler.PingHandler(zap.NewNop(), &statsd.Client{})(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
