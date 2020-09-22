package handler_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"urlshortner/pkg/http/internal/handler"
)

func TestPingHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/ping", nil)
	require.NoError(t, err)

	handler.PingHandler()(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"data\":\"pong\",\"error\":{},\"success\":true}", w.Body.String())
}
